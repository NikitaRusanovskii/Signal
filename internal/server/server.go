package server

import (
	"context"
	"net/http"
	"net/netip"
	"signal/internal/domain"
	"signal/internal/repository"
	"signal/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpHandler struct {
	Router  *gin.Engine
	Service *service.PeerService
}

func InitHandler(pr *repository.PeerManager) *HttpHandler {
	r := gin.Default()
	service := service.NewPeerService(pr)

	handler := HttpHandler{Router: r, Service: service}

	handler.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		addrPortRaw := c.Request.RemoteAddr
		addrPort := netip.MustParseAddrPort(addrPortRaw)
		handler.Service.UpdatePeerOnlineStatusByIP(c, true, addrPort)

		go func(addr netip.AddrPort) {
			time.Sleep(5 * time.Second)
			ctx := context.Background()
			handler.Service.UpdatePeerOnlineStatusByIP(ctx, false, addr)
		}(addrPort)
	})

	handler.Router.GET("/connect", func(c *gin.Context) {
		addrPortRaw := c.Request.RemoteAddr
		var role domain.Role = domain.SlaveRole
		tableEmpty := handler.Service.IsEmpty(c)
		if tableEmpty {
			role = domain.MasterRole
		}

		peer := domain.NewPeer(uuid.New(), role, true, addrPortRaw, time.Now())
		err := handler.Service.RegisterPeer(c, peer)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "connection error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "connected"})
	})

	handler.Router.GET("/master_ip", func(c *gin.Context) {
		masterIp, err := handler.Service.GetMasterIP(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "master ip getting error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"master_ip": masterIp.String()})
	})

	handler.Router.GET("/slave_ip", func(c *gin.Context) {
		slaveIp, err := handler.Service.GetLastSlavePeerIP(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "slave ip getting error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"slave_ip": slaveIp.String()})
	})

	return &handler
}

func (h *HttpHandler) Run() error {
	err := h.Router.Run()
	return err
}

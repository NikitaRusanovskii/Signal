package api

import (
	"net/http"
	"net/netip"
	"signal/internal/domain"
	"signal/internal/repository"
	"signal/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	service *service.PeerService
}

type CreatePeerRequest struct {
	Role domain.Role `json:"role" binding:"required"`
}

type SetMasterRequest struct {
	AddrPort string `json:"addr_port" binding:"required"`
}

func InitServer(pr *repository.PeerManager) *Server {
	r := gin.Default()
	s := service.NewPeerService(pr)

	r.GET("/get_slaves", func(c *gin.Context) {
		addrs, err := s.GetSlavesAddrPort(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error when get slaves addrPort"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successful", "slave_addrs": addrs})
	})

	r.GET("/get_masters", func(c *gin.Context) {
		addrs, err := s.GetMastersAddrPort(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error when get masters addrPort"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successful", "master_addrs": addrs})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/connect", func(c *gin.Context) {
		var req CreatePeerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		peer := domain.NewPeer(req.Role, true, c.Request.RemoteAddr, time.Now())
		s.AddPeer(c, peer)
		c.JSON(http.StatusCreated, gin.H{"status": "ok", "addr": c.Request.RemoteAddr})
	})

	r.DELETE("/disconnect", func(c *gin.Context) {
		err := s.DeletePeer(c, netip.MustParseAddrPort(c.Request.RemoteAddr))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error when delete peer"})
			return
		}
	})

	r.POST("/set_master", func(c *gin.Context) {
		var req SetMasterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err := s.SetRole(c, netip.MustParseAddrPort(req.AddrPort), domain.MasterRole)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error when set role"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	r.PATCH("/heartbeat", func(c *gin.Context) {
		err := s.SetOnline(c, netip.MustParseAddrPort(c.Request.RemoteAddr), true)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "error when heartbeat"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	return &Server{router: r, service: s}
}

func (h *Server) Run() error {
	err := h.router.Run()
	return err
}

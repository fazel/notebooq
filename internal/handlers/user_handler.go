package handlers

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/fazel/notebooq/internal/config"
	"github.com/fazel/notebooq/internal/repository"
	"github.com/fazel/notebooq/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	svc *service.UserService
	cfg *config.Config
}

func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	repo := repository.NewUserRepo(db)
	svc := service.NewUserService(repo, cfg.JWTSecret, cfg.AccessTokenExp)
	return &UserHandler{svc: svc, cfg: cfg}
}

// sendVerificationEmail
func sendVerificationEmail(to, code string, cfg *config.Config) error {
	subject := "Verify your Notebooq account"
	body := fmt.Sprintf("Hello!\n\nYour Notebooq verification code is: %s\n\nThank you!", code)

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTPEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPEmail, cfg.SMTPPassword)
	return d.DialAndSend(m)
}

func generateCode() string {
	const letters = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := service.GenerateCode() // ۶ رقمی
	u, err := h.svc.CreateUser(req.Username, req.Password, req.Email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// فرض می‌کنیم ایمیل موفق ارسال شد و کاربر کد را دریافت می‌کند
	u.VerifyCode = code
	h.svc.Update(u)

	c.JSON(http.StatusOK, gin.H{
		"message": "signup successful, please verify your email",
		"code":    code, // فقط برای تست؛ در واقعیت ایمیل می‌رود
	})
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.svc.GetByUsername(req.Username)
	if err != nil /*|| u.VerifyCode != req.Code*/ {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	u.IsVerified = true
	u.VerifyCode = ""
	h.svc.Update(u)

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.svc.GetByUsername(req.Username)
	if err != nil || !u.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials or email not verified"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.svc.GenerateJWT(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u, err := h.svc.GetByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"verified": u.IsVerified,
	})
}

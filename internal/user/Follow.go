package user

import (
	"net/http"
	"note/internal/models"
	"note/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *UserHandler) FollowUser(c *gin.Context) {
	targetIDstr := c.Param("id")
	targetID, err := strconv.ParseUint(targetIDstr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid target ID")
		return
	}

	me, err := utils.GetUserID(c)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	if me == uint(targetID) {
		utils.Error(c, http.StatusBadRequest, "You can't follow yourself")
		return // 记得 return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		followRel := models.UserFollow{
			FollowedID: uint(targetID),
			FollowerID: me,
		}

		if err := tx.Create(&followRel).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("id = ?", me).
			Update("follow_count", gorm.Expr("follow_count + 1")).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("id = ?", uint(targetID)).
			Update("fan_count", gorm.Expr("fan_count + 1")).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to follow: "+err.Error())
		return
	}
	utils.Success(c, gin.H{"message": "Followed successfully"})
}

func (h *UserHandler) UnfollowUser(c *gin.Context) {
	targetIDstr := c.Param("id")
	targetID, err := strconv.ParseUint(targetIDstr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid target ID")
		return
	}

	me, err := utils.GetUserID(c)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("follower_id = ? AND followed_id = ?", me, targetID).Delete(&models.UserFollow{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return nil
		}

		if err := tx.Model(&models.User{}).Where("id = ?", me).
			Update("follow_count", gorm.Expr("follow_count - 1")).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("id = ?", uint(targetID)).
			Update("fan_count", gorm.Expr("fan_count - 1")).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, gin.H{"message": "Unfollowed successfully"})
}

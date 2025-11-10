package user_s

import (
	"backend/logs"
	"backend/pkg/cryptoutil"
	"context"
	"errors"
)

func (s *UserServiceImpl) ResetPassword(ctx context.Context, token, newPassword string) error {
	email, err := s.emailRepo.GetEmailByResetPwdToken(ctx, token)
	if err != nil {
		logs.Sugar.Warnw("无效或过期的重置 token", "token", token)
		return errors.New("无效或过期的重置链接")
	}

	// 根据 email 查用户是否
	if exists, err := s.userRepo.ExistsByEmail(ctx, email); err != nil || !exists {
		logs.Sugar.Warnw("根据 email 查不到用户", "email", email)
		return errors.New("用户异常")
	}

	// 更新密码
	hashedPassword, err := cryptoutil.HashPassword(newPassword)
	if err != nil {
		logs.Sugar.Errorw("密码加密失败", "error", err)
		return errors.New("密码加密失败")
	}
	if err := s.userRepo.UpdatePasswordByEmail(ctx, email, hashedPassword); err != nil {
		logs.Sugar.Errorw("密码更新失败", "error", err)
		return errors.New("密码更新失败")
	}

	logs.Sugar.Infow("重置密码", "detail", "重置密码成功。")
	return nil
}

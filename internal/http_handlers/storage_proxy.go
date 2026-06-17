package http_handlers

import (
	"context"
	"fmt"

	"github.com/authorizerdev/authorizer/internal/graph/model"
	"github.com/authorizerdev/authorizer/internal/storage"
	"github.com/authorizerdev/authorizer/internal/storage/schemas"
	"github.com/authorizerdev/authorizer/internal/tenant"
)

// ContextBoundStorage delegates storage calls to the provider stored in request context.
type ContextBoundStorage struct{}

func (ContextBoundStorage) delegate(ctx context.Context) (storage.Provider, error) {
	return tenant.StorageFromContext(ctx)
}

func (c ContextBoundStorage) AddUser(ctx context.Context, user *schemas.User) (*schemas.User, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.AddUser(ctx, user)
}

func (c ContextBoundStorage) UpdateUser(ctx context.Context, user *schemas.User) (*schemas.User, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.UpdateUser(ctx, user)
}

func (c ContextBoundStorage) DeleteUser(ctx context.Context, user *schemas.User) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteUser(ctx, user)
}

func (c ContextBoundStorage) ListUsers(ctx context.Context, pagination *model.Pagination) ([]*schemas.User, *model.Pagination, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return p.ListUsers(ctx, pagination)
}

func (c ContextBoundStorage) GetUserByEmail(ctx context.Context, email string) (*schemas.User, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetUserByEmail(ctx, email)
}

func (c ContextBoundStorage) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*schemas.User, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetUserByPhoneNumber(ctx, phoneNumber)
}

func (c ContextBoundStorage) GetUserByID(ctx context.Context, id string) (*schemas.User, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetUserByID(ctx, id)
}

func (c ContextBoundStorage) UpdateUsers(ctx context.Context, data map[string]interface{}, ids []string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.UpdateUsers(ctx, data, ids)
}

func (c ContextBoundStorage) AddVerificationRequest(ctx context.Context, verificationRequest *schemas.VerificationRequest) (*schemas.VerificationRequest, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.AddVerificationRequest(ctx, verificationRequest)
}

func (c ContextBoundStorage) GetVerificationRequestByToken(ctx context.Context, token string) (*schemas.VerificationRequest, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetVerificationRequestByToken(ctx, token)
}

func (c ContextBoundStorage) GetVerificationRequestByEmail(ctx context.Context, email string, identifier string) (*schemas.VerificationRequest, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetVerificationRequestByEmail(ctx, email, identifier)
}

func (c ContextBoundStorage) ListVerificationRequests(ctx context.Context, pagination *model.Pagination) ([]*schemas.VerificationRequest, *model.Pagination, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return p.ListVerificationRequests(ctx, pagination)
}

func (c ContextBoundStorage) DeleteVerificationRequest(ctx context.Context, verificationRequest *schemas.VerificationRequest) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteVerificationRequest(ctx, verificationRequest)
}

func (c ContextBoundStorage) AddSession(ctx context.Context, session *schemas.Session) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.AddSession(ctx, session)
}

func (c ContextBoundStorage) DeleteSession(ctx context.Context, userId string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteSession(ctx, userId)
}

func (c ContextBoundStorage) AddWebhook(ctx context.Context, webhook *schemas.Webhook) (*schemas.Webhook, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.AddWebhook(ctx, webhook)
}

func (c ContextBoundStorage) UpdateWebhook(ctx context.Context, webhook *schemas.Webhook) (*schemas.Webhook, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.UpdateWebhook(ctx, webhook)
}

func (c ContextBoundStorage) ListWebhook(ctx context.Context, pagination *model.Pagination) ([]*schemas.Webhook, *model.Pagination, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return p.ListWebhook(ctx, pagination)
}

func (c ContextBoundStorage) GetWebhookByID(ctx context.Context, webhookID string) (*schemas.Webhook, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetWebhookByID(ctx, webhookID)
}

func (c ContextBoundStorage) GetWebhookByEventName(ctx context.Context, eventName string) ([]*schemas.Webhook, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetWebhookByEventName(ctx, eventName)
}

func (c ContextBoundStorage) DeleteWebhook(ctx context.Context, webhook *schemas.Webhook) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteWebhook(ctx, webhook)
}

func (c ContextBoundStorage) AddWebhookLog(ctx context.Context, webhookLog *schemas.WebhookLog) (*schemas.WebhookLog, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.AddWebhookLog(ctx, webhookLog)
}

func (c ContextBoundStorage) ListWebhookLogs(ctx context.Context, pagination *model.Pagination, webhookID string) ([]*schemas.WebhookLog, *model.Pagination, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return p.ListWebhookLogs(ctx, pagination, webhookID)
}

func (c ContextBoundStorage) AddEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) (*schemas.EmailTemplate, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.AddEmailTemplate(ctx, emailTemplate)
}

func (c ContextBoundStorage) UpdateEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) (*schemas.EmailTemplate, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.UpdateEmailTemplate(ctx, emailTemplate)
}

func (c ContextBoundStorage) ListEmailTemplate(ctx context.Context, pagination *model.Pagination) ([]*schemas.EmailTemplate, *model.Pagination, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return p.ListEmailTemplate(ctx, pagination)
}

func (c ContextBoundStorage) GetEmailTemplateByID(ctx context.Context, emailTemplateID string) (*schemas.EmailTemplate, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetEmailTemplateByID(ctx, emailTemplateID)
}

func (c ContextBoundStorage) GetEmailTemplateByEventName(ctx context.Context, eventName string) (*schemas.EmailTemplate, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetEmailTemplateByEventName(ctx, eventName)
}

func (c ContextBoundStorage) DeleteEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteEmailTemplate(ctx, emailTemplate)
}

func (c ContextBoundStorage) UpsertOTP(ctx context.Context, otp *schemas.OTP) (*schemas.OTP, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.UpsertOTP(ctx, otp)
}

func (c ContextBoundStorage) GetOTPByEmail(ctx context.Context, emailAddress string) (*schemas.OTP, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetOTPByEmail(ctx, emailAddress)
}

func (c ContextBoundStorage) GetOTPByPhoneNumber(ctx context.Context, phoneNumber string) (*schemas.OTP, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetOTPByPhoneNumber(ctx, phoneNumber)
}

func (c ContextBoundStorage) DeleteOTP(ctx context.Context, otp *schemas.OTP) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteOTP(ctx, otp)
}

func (c ContextBoundStorage) AddAuthenticator(ctx context.Context, totp *schemas.Authenticator) (*schemas.Authenticator, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.AddAuthenticator(ctx, totp)
}

func (c ContextBoundStorage) UpdateAuthenticator(ctx context.Context, totp *schemas.Authenticator) (*schemas.Authenticator, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.UpdateAuthenticator(ctx, totp)
}

func (c ContextBoundStorage) GetAuthenticatorDetailsByUserId(ctx context.Context, userId string, authenticatorType string) (*schemas.Authenticator, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetAuthenticatorDetailsByUserId(ctx, userId, authenticatorType)
}

func (c ContextBoundStorage) AddSessionToken(ctx context.Context, token *schemas.SessionToken) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.AddSessionToken(ctx, token)
}

func (c ContextBoundStorage) GetSessionTokenByUserIDAndKey(ctx context.Context, userId, key string) (*schemas.SessionToken, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetSessionTokenByUserIDAndKey(ctx, userId, key)
}

func (c ContextBoundStorage) DeleteSessionToken(ctx context.Context, id string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteSessionToken(ctx, id)
}

func (c ContextBoundStorage) DeleteSessionTokenByUserIDAndKey(ctx context.Context, userId, key string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteSessionTokenByUserIDAndKey(ctx, userId, key)
}

func (c ContextBoundStorage) DeleteAllSessionTokensByUserID(ctx context.Context, userId string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteAllSessionTokensByUserID(ctx, userId)
}

func (c ContextBoundStorage) DeleteSessionTokensByNamespace(ctx context.Context, namespace string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteSessionTokensByNamespace(ctx, namespace)
}

func (c ContextBoundStorage) CleanExpiredSessionTokens(ctx context.Context) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.CleanExpiredSessionTokens(ctx)
}

func (c ContextBoundStorage) GetAllSessionTokens(ctx context.Context) ([]*schemas.SessionToken, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetAllSessionTokens(ctx)
}

func (c ContextBoundStorage) AddMFASession(ctx context.Context, session *schemas.MFASession) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.AddMFASession(ctx, session)
}

func (c ContextBoundStorage) GetMFASessionByUserIDAndKey(ctx context.Context, userId, key string) (*schemas.MFASession, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetMFASessionByUserIDAndKey(ctx, userId, key)
}

func (c ContextBoundStorage) DeleteMFASession(ctx context.Context, id string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteMFASession(ctx, id)
}

func (c ContextBoundStorage) DeleteMFASessionByUserIDAndKey(ctx context.Context, userId, key string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteMFASessionByUserIDAndKey(ctx, userId, key)
}

func (c ContextBoundStorage) GetAllMFASessionsByUserID(ctx context.Context, userId string) ([]*schemas.MFASession, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetAllMFASessionsByUserID(ctx, userId)
}

func (c ContextBoundStorage) CleanExpiredMFASessions(ctx context.Context) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.CleanExpiredMFASessions(ctx)
}

func (c ContextBoundStorage) GetAllMFASessions(ctx context.Context) ([]*schemas.MFASession, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetAllMFASessions(ctx)
}

func (c ContextBoundStorage) AddOAuthState(ctx context.Context, state *schemas.OAuthState) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.AddOAuthState(ctx, state)
}

func (c ContextBoundStorage) GetOAuthStateByKey(ctx context.Context, key string) (*schemas.OAuthState, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetOAuthStateByKey(ctx, key)
}

func (c ContextBoundStorage) DeleteOAuthStateByKey(ctx context.Context, key string) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteOAuthStateByKey(ctx, key)
}

func (c ContextBoundStorage) GetAllOAuthStates(ctx context.Context) ([]*schemas.OAuthState, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, err
	}
	return p.GetAllOAuthStates(ctx)
}

func (c ContextBoundStorage) AddAuditLog(ctx context.Context, log *schemas.AuditLog) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.AddAuditLog(ctx, log)
}

func (c ContextBoundStorage) ListAuditLogs(ctx context.Context, pagination *model.Pagination, filter map[string]interface{}) ([]*schemas.AuditLog, *model.Pagination, error) {
	p, err := c.delegate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return p.ListAuditLogs(ctx, pagination, filter)
}

func (c ContextBoundStorage) DeleteAuditLogsBefore(ctx context.Context, before int64) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return err
	}
	return p.DeleteAuditLogsBefore(ctx, before)
}

func (c ContextBoundStorage) HealthCheck(ctx context.Context) error {
	p, err := c.delegate(ctx)
	if err != nil {
		return fmt.Errorf("health check requires tenant context: %w", err)
	}
	return p.HealthCheck(ctx)
}

func (c ContextBoundStorage) Close() error {
	return nil
}

var _ storage.Provider = ContextBoundStorage{}

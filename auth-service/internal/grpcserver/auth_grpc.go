package grpcserver

import (
	"context"

	"neocentral-go/auth-service/internal/domain"
	"neocentral-go/auth-service/internal/repository"
	"neocentral-go/auth-service/internal/service"
	pb "neocentral-go/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthGRPCServer implements the auth.AuthServiceServer gRPC interface.
type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	authSvc  *service.AuthService
	userRepo repository.UserRepository
}

func NewAuthGRPCServer(authSvc *service.AuthService, userRepo repository.UserRepository) *AuthGRPCServer {
	return &AuthGRPCServer{
		authSvc:  authSvc,
		userRepo: userRepo,
	}
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := s.authSvc.VerifyAccessToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	user, err := s.userRepo.FindByID(ctx, claims.Subject)
	if err != nil || user == nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	var roleNames []string
	for _, uhr := range user.UserHasRoles {
		roleNames = append(roleNames, uhr.Role.Name)
	}

	return &pb.ValidateTokenResponse{
		Valid:          true,
		UserId:         user.ID,
		IdentityNumber: user.IdentityNumber,
		FullName:       user.FullName,
		Email:          ptrStr(user.Email),
		Roles:          roleNames,
	}, nil
}

func (s *AuthGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", req.UserId)
	}
	return toUserProto(user), nil
}

func (s *AuthGRPCServer) GetUserContact(ctx context.Context, req *pb.GetUserContactRequest) (*pb.UserContactResponse, error) {
	user, err := s.userRepo.FindByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &pb.UserContactResponse{
		UserId:      user.ID,
		FullName:    user.FullName,
		Email:       ptrStr(user.Email),
		PhoneNumber: ptrStr(user.PhoneNumber),
	}, nil
}

func (s *AuthGRPCServer) HasRole(ctx context.Context, req *pb.HasRoleRequest) (*pb.HasRoleResponse, error) {
	has, err := s.userRepo.HasRole(ctx, req.UserId, req.RoleName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check role: %v", err)
	}
	return &pb.HasRoleResponse{HasRole: has}, nil
}

func (s *AuthGRPCServer) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	users, err := s.userRepo.BatchGetByIDs(ctx, req.UserIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to batch get users: %v", err)
	}
	var results []*pb.UserResponse
	for i := range users {
		results = append(results, toUserProto(&users[i]))
	}
	return &pb.BatchGetUsersResponse{Users: results}, nil
}

func (s *AuthGRPCServer) GetUsersByRole(ctx context.Context, req *pb.GetUsersByRoleRequest) (*pb.GetUsersByRoleResponse, error) {
	page, pageSize := int(req.Page), int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	users, total, err := s.userRepo.GetUsersByRole(ctx, req.RoleName, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users by role: %v", err)
	}

	var results []*pb.UserResponse
	for i := range users {
		results = append(results, toUserProto(&users[i]))
	}
	return &pb.GetUsersByRoleResponse{Users: results, Total: int32(total)}, nil
}

// ── Helpers ──────────────────────────────────────────────────────

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toUserProto(u *domain.User) *pb.UserResponse {
	resp := &pb.UserResponse{
		Id:             u.ID,
		FullName:       u.FullName,
		IdentityNumber: u.IdentityNumber,
		IdentityType:   string(u.IdentityType),
		IsVerified:     u.IsVerified,
		Email:          ptrStr(u.Email),
		PhoneNumber:    ptrStr(u.PhoneNumber),
		AvatarUrl:      ptrStr(u.AvatarURL),
	}

	for _, uhr := range u.UserHasRoles {
		resp.Roles = append(resp.Roles, &pb.RoleInfo{
			Id:     uhr.RoleID,
			Name:   uhr.Role.Name,
			Status: string(uhr.Status),
		})
	}

	if u.Student != nil {
		resp.Student = &pb.StudentInfo{
			SksCompleted: int32(u.Student.SKSCompleted),
			Status:       string(u.Student.Status),
		}
		if u.Student.EnrollmentYear != nil {
			resp.Student.EnrollmentYear = int32(*u.Student.EnrollmentYear)
		}
		if u.Student.CurrentSemester != nil {
			resp.Student.CurrentSemester = int32(*u.Student.CurrentSemester)
		}
	}

	if u.Lecturer != nil {
		resp.Lecturer = &pb.LecturerInfo{
			ScienceGroupId: ptrStr(u.Lecturer.ScienceGroupID),
		}
	}

	return resp
}

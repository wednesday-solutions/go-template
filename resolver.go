// Package goboiler ...
// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
package goboiler

import (
	"context"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/gqlgen/helper"
	fm "github.com/wednesday-solutions/go-boiler/graphql_models"
	hm "github.com/wednesday-solutions/go-boiler/helpers"
	dm "github.com/wednesday-solutions/go-boiler/models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/middleware/auth"
)

// Resolver ...
type Resolver struct {
}

const inputKey = "input"

func (r *mutationResolver) CreateComment(ctx context.Context, input fm.CommentCreateInput) (*fm.CommentPayload, error) {

	m := hm.CommentCreateInputToBoiler(&input)

	m.UserID = auth.UserIDFromContext(ctx)

	whiteList := hm.CommentCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
		dm.CommentColumns.UserID,
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.CommentPreloadMap, hm.CommentPayloadPreloadLevels.Comment)
	mods = append(mods, dm.CommentWhere.ID.EQ(m.ID))
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	pM, err := dm.Comments(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CommentPayload{
		Comment: hm.CommentToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateComments(ctx context.Context, input fm.CommentsCreateInput) (*fm.CommentsPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, input fm.CommentUpdateInput) (*fm.CommentPayload, error) {
	m := hm.CommentUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.CommentID(id)
	if _, err := dm.Comments(
		dm.CommentWhere.ID.EQ(dbID),
		dm.CommentWhere.UserID.EQ(auth.UserIDFromContext(ctx)),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.CommentPreloadMap, hm.CommentPayloadPreloadLevels.Comment)
	mods = append(mods, dm.CommentWhere.ID.EQ(dbID))
	mods = append(mods, dm.CommentWhere.UserID.EQ(auth.UserIDFromContext(ctx)))

	pM, err := dm.Comments(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CommentPayload{
		Comment: hm.CommentToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateComments(ctx context.Context, filter *fm.CommentFilter, input fm.CommentUpdateInput) (*fm.CommentsUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, hm.CommentFilterToMods(filter)...)

	m := hm.CommentUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Comments(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.CommentsUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*fm.CommentDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.CommentWhere.ID.EQ(hm.CommentID(id)),
		dm.CommentWhere.UserID.EQ(auth.UserIDFromContext(ctx)),
	}
	_, err := dm.Comments(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.CommentDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteComments(ctx context.Context, filter *fm.CommentFilter) (*fm.CommentsDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, hm.CommentFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.CommentColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Comments))

	var IDsToRemove []helper.RemovedID
	if err := dm.Comments(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Comments(dm.CommentWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.CommentsDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Comments),
	}, nil
}

func (r *mutationResolver) CreateCompany(ctx context.Context, input fm.CompanyCreateInput) (*fm.CompanyPayload, error) {

	m := hm.CompanyCreateInputToBoiler(&input)

	whiteList := hm.CompanyCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.CompanyPreloadMap, hm.CompanyPayloadPreloadLevels.Company)
	mods = append(mods, dm.CompanyWhere.ID.EQ(m.ID))
	pM, err := dm.Companies(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CompanyPayload{
		Company: hm.CompanyToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateCompanies(ctx context.Context, input fm.CompaniesCreateInput) (*fm.CompaniesPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateCompany(ctx context.Context, id string, input fm.CompanyUpdateInput) (*fm.CompanyPayload, error) {
	m := hm.CompanyUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.CompanyID(id)
	if _, err := dm.Companies(
		dm.CompanyWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.CompanyPreloadMap, hm.CompanyPayloadPreloadLevels.Company)
	mods = append(mods, dm.CompanyWhere.ID.EQ(dbID))

	pM, err := dm.Companies(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CompanyPayload{
		Company: hm.CompanyToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateCompanies(ctx context.Context, filter *fm.CompanyFilter, input fm.CompanyUpdateInput) (*fm.CompaniesUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.CompanyFilterToMods(filter)...)

	m := hm.CompanyUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Companies(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.CompaniesUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteCompany(ctx context.Context, id string) (*fm.CompanyDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.CompanyWhere.ID.EQ(hm.CompanyID(id)),
	}
	_, err := dm.Companies(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.CompanyDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteCompanies(ctx context.Context, filter *fm.CompanyFilter) (*fm.CompaniesDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.CompanyFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.CompanyColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Companies))

	var IDsToRemove []helper.RemovedID
	if err := dm.Companies(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Companies(dm.CompanyWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.CompaniesDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Companies),
	}, nil
}

func (r *mutationResolver) CreateFollower(ctx context.Context, input fm.FollowerCreateInput) (*fm.FollowerPayload, error) {

	m := hm.FollowerCreateInputToBoiler(&input)

	whiteList := hm.FollowerCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.FollowerPreloadMap, hm.FollowerPayloadPreloadLevels.Follower)
	mods = append(mods, dm.FollowerWhere.ID.EQ(m.ID))
	pM, err := dm.Followers(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.FollowerPayload{
		Follower: hm.FollowerToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateFollowers(ctx context.Context, input fm.FollowersCreateInput) (*fm.FollowersPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateFollower(ctx context.Context, id string, input fm.FollowerUpdateInput) (*fm.FollowerPayload, error) {
	m := hm.FollowerUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.FollowerID(id)
	if _, err := dm.Followers(
		dm.FollowerWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.FollowerPreloadMap, hm.FollowerPayloadPreloadLevels.Follower)
	mods = append(mods, dm.FollowerWhere.ID.EQ(dbID))

	pM, err := dm.Followers(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.FollowerPayload{
		Follower: hm.FollowerToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateFollowers(ctx context.Context, filter *fm.FollowerFilter, input fm.FollowerUpdateInput) (*fm.FollowersUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.FollowerFilterToMods(filter)...)

	m := hm.FollowerUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Followers(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.FollowersUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteFollower(ctx context.Context, id string) (*fm.FollowerDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.FollowerWhere.ID.EQ(hm.FollowerID(id)),
	}
	_, err := dm.Followers(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.FollowerDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteFollowers(ctx context.Context, filter *fm.FollowerFilter) (*fm.FollowersDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.FollowerFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.FollowerColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Followers))

	var IDsToRemove []helper.RemovedID
	if err := dm.Followers(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Followers(dm.FollowerWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.FollowersDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Followers),
	}, nil
}

func (r *mutationResolver) CreateLocation(ctx context.Context, input fm.LocationCreateInput) (*fm.LocationPayload, error) {

	m := hm.LocationCreateInputToBoiler(&input)

	whiteList := hm.LocationCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.LocationPreloadMap, hm.LocationPayloadPreloadLevels.Location)
	mods = append(mods, dm.LocationWhere.ID.EQ(m.ID))
	pM, err := dm.Locations(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.LocationPayload{
		Location: hm.LocationToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateLocations(ctx context.Context, input fm.LocationsCreateInput) (*fm.LocationsPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input fm.LocationUpdateInput) (*fm.LocationPayload, error) {
	m := hm.LocationUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.LocationID(id)
	if _, err := dm.Locations(
		dm.LocationWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.LocationPreloadMap, hm.LocationPayloadPreloadLevels.Location)
	mods = append(mods, dm.LocationWhere.ID.EQ(dbID))

	pM, err := dm.Locations(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.LocationPayload{
		Location: hm.LocationToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateLocations(ctx context.Context, filter *fm.LocationFilter, input fm.LocationUpdateInput) (*fm.LocationsUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.LocationFilterToMods(filter)...)

	m := hm.LocationUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Locations(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.LocationsUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteLocation(ctx context.Context, id string) (*fm.LocationDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.LocationWhere.ID.EQ(hm.LocationID(id)),
	}
	_, err := dm.Locations(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.LocationDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteLocations(ctx context.Context, filter *fm.LocationFilter) (*fm.LocationsDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.LocationFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.LocationColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Locations))

	var IDsToRemove []helper.RemovedID
	if err := dm.Locations(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Locations(dm.LocationWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.LocationsDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Locations),
	}, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input fm.PostCreateInput) (*fm.PostPayload, error) {

	m := hm.PostCreateInputToBoiler(&input)

	m.UserID = auth.UserIDFromContext(ctx)

	whiteList := hm.PostCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
		dm.PostColumns.UserID,
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.PostPreloadMap, hm.PostPayloadPreloadLevels.Post)
	mods = append(mods, dm.PostWhere.ID.EQ(m.ID))
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	pM, err := dm.Posts(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.PostPayload{
		Post: hm.PostToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreatePosts(ctx context.Context, input fm.PostsCreateInput) (*fm.PostsPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input fm.PostUpdateInput) (*fm.PostPayload, error) {
	m := hm.PostUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.PostID(id)
	if _, err := dm.Posts(
		dm.PostWhere.ID.EQ(dbID),
		dm.PostWhere.UserID.EQ(auth.UserIDFromContext(ctx)),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.PostPreloadMap, hm.PostPayloadPreloadLevels.Post)
	mods = append(mods, dm.PostWhere.ID.EQ(dbID))
	mods = append(mods, dm.PostWhere.UserID.EQ(auth.UserIDFromContext(ctx)))

	pM, err := dm.Posts(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.PostPayload{
		Post: hm.PostToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdatePosts(ctx context.Context, filter *fm.PostFilter, input fm.PostUpdateInput) (*fm.PostsUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, hm.PostFilterToMods(filter)...)

	m := hm.PostUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Posts(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.PostsUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*fm.PostDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.PostWhere.ID.EQ(hm.PostID(id)),
		dm.PostWhere.UserID.EQ(auth.UserIDFromContext(ctx)),
	}
	_, err := dm.Posts(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.PostDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeletePosts(ctx context.Context, filter *fm.PostFilter) (*fm.PostsDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, hm.PostFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.PostColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Posts))

	var IDsToRemove []helper.RemovedID
	if err := dm.Posts(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Posts(dm.PostWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.PostsDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Posts),
	}, nil
}

func (r *mutationResolver) CreateRole(ctx context.Context, input fm.RoleCreateInput) (*fm.RolePayload, error) {

	m := hm.RoleCreateInputToBoiler(&input)

	whiteList := hm.RoleCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.RolePreloadMap, hm.RolePayloadPreloadLevels.Role)
	mods = append(mods, dm.RoleWhere.ID.EQ(m.ID))
	pM, err := dm.Roles(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.RolePayload{
		Role: hm.RoleToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateRoles(ctx context.Context, input fm.RolesCreateInput) (*fm.RolesPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateRole(ctx context.Context, id string, input fm.RoleUpdateInput) (*fm.RolePayload, error) {
	m := hm.RoleUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.RoleID(id)
	if _, err := dm.Roles(
		dm.RoleWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.RolePreloadMap, hm.RolePayloadPreloadLevels.Role)
	mods = append(mods, dm.RoleWhere.ID.EQ(dbID))

	pM, err := dm.Roles(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.RolePayload{
		Role: hm.RoleToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateRoles(ctx context.Context, filter *fm.RoleFilter, input fm.RoleUpdateInput) (*fm.RolesUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.RoleFilterToMods(filter)...)

	m := hm.RoleUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Roles(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.RolesUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteRole(ctx context.Context, id string) (*fm.RoleDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.RoleWhere.ID.EQ(hm.RoleID(id)),
	}
	_, err := dm.Roles(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.RoleDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteRoles(ctx context.Context, filter *fm.RoleFilter) (*fm.RolesDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.RoleFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.RoleColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Roles))

	var IDsToRemove []helper.RemovedID
	if err := dm.Roles(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Roles(dm.RoleWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.RolesDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Roles),
	}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input fm.UserCreateInput) (*fm.UserPayload, error) {

	m := hm.UserCreateInputToBoiler(&input)

	whiteList := hm.UserCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.UserPreloadMap, hm.UserPayloadPreloadLevels.User)
	mods = append(mods, dm.UserWhere.ID.EQ(m.ID))
	pM, err := dm.Users(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.UserPayload{
		User: hm.UserToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateUsers(ctx context.Context, input fm.UsersCreateInput) (*fm.UsersPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input fm.UserUpdateInput) (*fm.UserPayload, error) {
	m := hm.UserUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := hm.UserID(id)
	if _, err := dm.Users(
		dm.UserWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, hm.UserPreloadMap, hm.UserPayloadPreloadLevels.User)
	mods = append(mods, dm.UserWhere.ID.EQ(dbID))

	pM, err := dm.Users(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.UserPayload{
		User: hm.UserToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateUsers(ctx context.Context, filter *fm.UserFilter, input fm.UserUpdateInput) (*fm.UsersUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.UserFilterToMods(filter)...)

	m := hm.UserUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Users(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.UsersUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*fm.UserDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.UserWhere.ID.EQ(hm.UserID(id)),
	}
	_, err := dm.Users(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.UserDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, filter *fm.UserFilter) (*fm.UsersDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, hm.UserFilterToMods(filter)...)
	mods = append(mods, qm.Select(dm.UserColumns.ID))
	mods = append(mods, qm.From(dm.TableNames.Users))

	var IDsToRemove []helper.RemovedID
	if err := dm.Users(mods...).Bind(context.Background(), boil.GetContextDB(), IDsToRemove); err != nil {
		return nil, err
	}

	boilerIDs := helper.RemovedIDsToInt(IDsToRemove)
	if _, err := dm.Users(dm.UserWhere.ID.IN(boilerIDs)).DeleteAll(context.Background(), boil.GetContextDB()); err != nil {
		return nil, err
	}

	return &fm.UsersDeletePayload{
		Ids: helper.IDsToGraphQL(boilerIDs, dm.TableNames.Users),
	}, nil
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*fm.Comment, error) {
	dbID := hm.CommentID(id)
	mods := helper.GetPreloadMods(ctx, hm.CommentPreloadMap)
	mods = append(mods, dm.CommentWhere.ID.EQ(dbID))
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	m, err := dm.Comments(mods...).One(context.Background(), boil.GetContextDB())
	return hm.CommentToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Comments(ctx context.Context, filter *fm.CommentFilter, pagination *fm.CommentPagination) ([]*fm.Comment, error) {
	mods := helper.GetPreloadMods(ctx, hm.CommentPreloadMap)
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, hm.CommentPaginationToMods(pagination)...)
	mods = append(mods, hm.CommentFilterToMods(filter)...)
	a, err := dm.Comments(mods...).All(context.Background(), boil.GetContextDB())
	return hm.CommentsToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Company(ctx context.Context, id string) (*fm.Company, error) {
	dbID := hm.CompanyID(id)
	mods := helper.GetPreloadMods(ctx, hm.CompanyPreloadMap)
	mods = append(mods, dm.CompanyWhere.ID.EQ(dbID))
	m, err := dm.Companies(mods...).One(context.Background(), boil.GetContextDB())
	return hm.CompanyToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Companies(ctx context.Context, filter *fm.CompanyFilter, pagination *fm.CompanyPagination) ([]*fm.Company, error) {
	mods := helper.GetPreloadMods(ctx, hm.CompanyPreloadMap)
	mods = append(mods, hm.CompanyPaginationToMods(pagination)...)
	mods = append(mods, hm.CompanyFilterToMods(filter)...)
	a, err := dm.Companies(mods...).All(context.Background(), boil.GetContextDB())
	return hm.CompaniesToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Follower(ctx context.Context, id string) (*fm.Follower, error) {
	dbID := hm.FollowerID(id)
	mods := helper.GetPreloadMods(ctx, hm.FollowerPreloadMap)
	mods = append(mods, dm.FollowerWhere.ID.EQ(dbID))
	m, err := dm.Followers(mods...).One(context.Background(), boil.GetContextDB())
	return hm.FollowerToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Followers(ctx context.Context, filter *fm.FollowerFilter, pagination *fm.FollowerPagination) ([]*fm.Follower, error) {
	mods := helper.GetPreloadMods(ctx, hm.FollowerPreloadMap)
	mods = append(mods, hm.FollowerPaginationToMods(pagination)...)
	mods = append(mods, hm.FollowerFilterToMods(filter)...)
	a, err := dm.Followers(mods...).All(context.Background(), boil.GetContextDB())
	return hm.FollowersToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Location(ctx context.Context, id string) (*fm.Location, error) {
	dbID := hm.LocationID(id)
	mods := helper.GetPreloadMods(ctx, hm.LocationPreloadMap)
	mods = append(mods, dm.LocationWhere.ID.EQ(dbID))
	m, err := dm.Locations(mods...).One(context.Background(), boil.GetContextDB())
	return hm.LocationToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Locations(ctx context.Context, filter *fm.LocationFilter, pagination *fm.LocationPagination) ([]*fm.Location, error) {
	mods := helper.GetPreloadMods(ctx, hm.LocationPreloadMap)
	mods = append(mods, hm.LocationPaginationToMods(pagination)...)
	mods = append(mods, hm.LocationFilterToMods(filter)...)
	a, err := dm.Locations(mods...).All(context.Background(), boil.GetContextDB())
	return hm.LocationsToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Post(ctx context.Context, id string) (*fm.Post, error) {
	dbID := hm.PostID(id)
	mods := helper.GetPreloadMods(ctx, hm.PostPreloadMap)
	mods = append(mods, dm.PostWhere.ID.EQ(dbID))
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	m, err := dm.Posts(mods...).One(context.Background(), boil.GetContextDB())
	return hm.PostToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Posts(ctx context.Context, filter *fm.PostFilter, pagination *fm.PostPagination) ([]*fm.Post, error) {
	mods := helper.GetPreloadMods(ctx, hm.PostPreloadMap)
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, hm.PostPaginationToMods(pagination)...)
	mods = append(mods, hm.PostFilterToMods(filter)...)
	a, err := dm.Posts(mods...).All(context.Background(), boil.GetContextDB())
	return hm.PostsToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Role(ctx context.Context, id string) (*fm.Role, error) {
	dbID := hm.RoleID(id)
	mods := helper.GetPreloadMods(ctx, hm.RolePreloadMap)
	mods = append(mods, dm.RoleWhere.ID.EQ(dbID))
	m, err := dm.Roles(mods...).One(context.Background(), boil.GetContextDB())
	return hm.RoleToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Roles(ctx context.Context, filter *fm.RoleFilter, pagination *fm.RolePagination) ([]*fm.Role, error) {
	mods := helper.GetPreloadMods(ctx, hm.RolePreloadMap)
	mods = append(mods, hm.RolePaginationToMods(pagination)...)
	mods = append(mods, hm.RoleFilterToMods(filter)...)
	a, err := dm.Roles(mods...).All(context.Background(), boil.GetContextDB())
	return hm.RolesToGraphQL(a, nil, 0), err
}

func (r *queryResolver) User(ctx context.Context, id string) (*fm.User, error) {
	dbID := hm.UserID(id)
	mods := helper.GetPreloadMods(ctx, hm.UserPreloadMap)
	mods = append(mods, dm.UserWhere.ID.EQ(dbID))
	m, err := dm.Users(mods...).One(context.Background(), boil.GetContextDB())
	return hm.UserToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Users(ctx context.Context, filter *fm.UserFilter, pagination *fm.UserPagination) ([]*fm.User, error) {
	mods := helper.GetPreloadMods(ctx, hm.UserPreloadMap)
	mods = append(mods, hm.UserPaginationToMods(pagination)...)
	mods = append(mods, hm.UserFilterToMods(filter)...)
	a, err := dm.Users(mods...).All(context.Background(), boil.GetContextDB())
	return hm.UsersToGraphQL(a, nil, 0), err
}

// Mutation ...
func (r *Resolver) Mutation() fm.MutationResolver { return &mutationResolver{r} }

// Query ...
func (r *Resolver) Query() fm.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

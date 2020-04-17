// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
package goboiler

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/wednesday-solutions/go-boiler/gqlgen/helper"
	fm "github.com/wednesday-solutions/go-boiler/graphql_models"
	. "github.com/wednesday-solutions/go-boiler/helpers"
	dm "github.com/wednesday-solutions/go-boiler/models"
	"github.com/wednesday-solutions/go-boiler/pkg/utl/middleware/auth"
)

type Resolver struct {
	db *sql.DB
}

const inputKey = "input"

func (r *mutationResolver) CreateComment(ctx context.Context, input fm.CommentCreateInput) (*fm.CommentPayload, error) {

	m := CommentCreateInputToBoiler(&input)

	m.UserID = auth.UserIDFromContext(ctx)

	whiteList := CommentCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
		dm.CommentColumns.UserID,
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, CommentPreloadMap, CommentPayloadPreloadLevels.Comment)
	mods = append(mods, dm.CommentWhere.ID.EQ(m.ID))
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	pM, err := dm.Comments(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CommentPayload{
		Comment: CommentToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateComments(ctx context.Context, input fm.CommentsCreateInput) (*fm.CommentsPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, input fm.CommentUpdateInput) (*fm.CommentPayload, error) {
	m := CommentUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := CommentID(id)
	if _, err := dm.Comments(
		dm.CommentWhere.ID.EQ(dbID),
		dm.CommentWhere.UserID.EQ(auth.UserIDFromContext(ctx)),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, CommentPreloadMap, CommentPayloadPreloadLevels.Comment)
	mods = append(mods, dm.CommentWhere.ID.EQ(dbID))
	mods = append(mods, dm.CommentWhere.UserID.EQ(auth.UserIDFromContext(ctx)))

	pM, err := dm.Comments(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CommentPayload{
		Comment: CommentToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateComments(ctx context.Context, filter *fm.CommentFilter, input fm.CommentUpdateInput) (*fm.CommentsUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, CommentFilterToMods(filter)...)

	m := CommentUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Comments(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.CommentsUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*fm.CommentDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.CommentWhere.ID.EQ(CommentID(id)),
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
	mods = append(mods, CommentFilterToMods(filter)...)
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

	m := CompanyCreateInputToBoiler(&input)

	whiteList := CompanyCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, CompanyPreloadMap, CompanyPayloadPreloadLevels.Company)
	mods = append(mods, dm.CompanyWhere.ID.EQ(m.ID))
	pM, err := dm.Companies(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CompanyPayload{
		Company: CompanyToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateCompanies(ctx context.Context, input fm.CompaniesCreateInput) (*fm.CompaniesPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateCompany(ctx context.Context, id string, input fm.CompanyUpdateInput) (*fm.CompanyPayload, error) {
	m := CompanyUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := CompanyID(id)
	if _, err := dm.Companies(
		dm.CompanyWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, CompanyPreloadMap, CompanyPayloadPreloadLevels.Company)
	mods = append(mods, dm.CompanyWhere.ID.EQ(dbID))

	pM, err := dm.Companies(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.CompanyPayload{
		Company: CompanyToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateCompanies(ctx context.Context, filter *fm.CompanyFilter, input fm.CompanyUpdateInput) (*fm.CompaniesUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, CompanyFilterToMods(filter)...)

	m := CompanyUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Companies(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.CompaniesUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteCompany(ctx context.Context, id string) (*fm.CompanyDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.CompanyWhere.ID.EQ(CompanyID(id)),
	}
	_, err := dm.Companies(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.CompanyDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteCompanies(ctx context.Context, filter *fm.CompanyFilter) (*fm.CompaniesDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, CompanyFilterToMods(filter)...)
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

	m := FollowerCreateInputToBoiler(&input)

	whiteList := FollowerCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, FollowerPreloadMap, FollowerPayloadPreloadLevels.Follower)
	mods = append(mods, dm.FollowerWhere.ID.EQ(m.ID))
	pM, err := dm.Followers(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.FollowerPayload{
		Follower: FollowerToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateFollowers(ctx context.Context, input fm.FollowersCreateInput) (*fm.FollowersPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateFollower(ctx context.Context, id string, input fm.FollowerUpdateInput) (*fm.FollowerPayload, error) {
	m := FollowerUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := FollowerID(id)
	if _, err := dm.Followers(
		dm.FollowerWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, FollowerPreloadMap, FollowerPayloadPreloadLevels.Follower)
	mods = append(mods, dm.FollowerWhere.ID.EQ(dbID))

	pM, err := dm.Followers(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.FollowerPayload{
		Follower: FollowerToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateFollowers(ctx context.Context, filter *fm.FollowerFilter, input fm.FollowerUpdateInput) (*fm.FollowersUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, FollowerFilterToMods(filter)...)

	m := FollowerUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Followers(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.FollowersUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteFollower(ctx context.Context, id string) (*fm.FollowerDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.FollowerWhere.ID.EQ(FollowerID(id)),
	}
	_, err := dm.Followers(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.FollowerDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteFollowers(ctx context.Context, filter *fm.FollowerFilter) (*fm.FollowersDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, FollowerFilterToMods(filter)...)
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

	m := LocationCreateInputToBoiler(&input)

	whiteList := LocationCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, LocationPreloadMap, LocationPayloadPreloadLevels.Location)
	mods = append(mods, dm.LocationWhere.ID.EQ(m.ID))
	pM, err := dm.Locations(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.LocationPayload{
		Location: LocationToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateLocations(ctx context.Context, input fm.LocationsCreateInput) (*fm.LocationsPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input fm.LocationUpdateInput) (*fm.LocationPayload, error) {
	m := LocationUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := LocationID(id)
	if _, err := dm.Locations(
		dm.LocationWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, LocationPreloadMap, LocationPayloadPreloadLevels.Location)
	mods = append(mods, dm.LocationWhere.ID.EQ(dbID))

	pM, err := dm.Locations(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.LocationPayload{
		Location: LocationToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateLocations(ctx context.Context, filter *fm.LocationFilter, input fm.LocationUpdateInput) (*fm.LocationsUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, LocationFilterToMods(filter)...)

	m := LocationUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Locations(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.LocationsUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteLocation(ctx context.Context, id string) (*fm.LocationDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.LocationWhere.ID.EQ(LocationID(id)),
	}
	_, err := dm.Locations(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.LocationDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteLocations(ctx context.Context, filter *fm.LocationFilter) (*fm.LocationsDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, LocationFilterToMods(filter)...)
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

	m := PostCreateInputToBoiler(&input)

	m.UserID = auth.UserIDFromContext(ctx)

	whiteList := PostCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
		dm.PostColumns.UserID,
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, PostPreloadMap, PostPayloadPreloadLevels.Post)
	mods = append(mods, dm.PostWhere.ID.EQ(m.ID))
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	pM, err := dm.Posts(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.PostPayload{
		Post: PostToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreatePosts(ctx context.Context, input fm.PostsCreateInput) (*fm.PostsPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input fm.PostUpdateInput) (*fm.PostPayload, error) {
	m := PostUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := PostID(id)
	if _, err := dm.Posts(
		dm.PostWhere.ID.EQ(dbID),
		dm.PostWhere.UserID.EQ(auth.UserIDFromContext(ctx)),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, PostPreloadMap, PostPayloadPreloadLevels.Post)
	mods = append(mods, dm.PostWhere.ID.EQ(dbID))
	mods = append(mods, dm.PostWhere.UserID.EQ(auth.UserIDFromContext(ctx)))

	pM, err := dm.Posts(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.PostPayload{
		Post: PostToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdatePosts(ctx context.Context, filter *fm.PostFilter, input fm.PostUpdateInput) (*fm.PostsUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, PostFilterToMods(filter)...)

	m := PostUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Posts(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.PostsUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*fm.PostDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.PostWhere.ID.EQ(PostID(id)),
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
	mods = append(mods, PostFilterToMods(filter)...)
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

	m := RoleCreateInputToBoiler(&input)

	whiteList := RoleCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, RolePreloadMap, RolePayloadPreloadLevels.Role)
	mods = append(mods, dm.RoleWhere.ID.EQ(m.ID))
	pM, err := dm.Roles(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.RolePayload{
		Role: RoleToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateRoles(ctx context.Context, input fm.RolesCreateInput) (*fm.RolesPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateRole(ctx context.Context, id string, input fm.RoleUpdateInput) (*fm.RolePayload, error) {
	m := RoleUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := RoleID(id)
	if _, err := dm.Roles(
		dm.RoleWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, RolePreloadMap, RolePayloadPreloadLevels.Role)
	mods = append(mods, dm.RoleWhere.ID.EQ(dbID))

	pM, err := dm.Roles(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.RolePayload{
		Role: RoleToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateRoles(ctx context.Context, filter *fm.RoleFilter, input fm.RoleUpdateInput) (*fm.RolesUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, RoleFilterToMods(filter)...)

	m := RoleUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Roles(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.RolesUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteRole(ctx context.Context, id string) (*fm.RoleDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.RoleWhere.ID.EQ(RoleID(id)),
	}
	_, err := dm.Roles(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.RoleDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteRoles(ctx context.Context, filter *fm.RoleFilter) (*fm.RolesDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, RoleFilterToMods(filter)...)
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

	m := UserCreateInputToBoiler(&input)

	whiteList := UserCreateInputToBoilerWhitelist(
		helper.GetInputFromContext(ctx, inputKey),
	)
	if err := m.Insert(context.Background(), boil.GetContextDB(), whiteList); err != nil {
		return nil, err
	}

	// resolve requested fields after creating
	mods := helper.GetPreloadModsWithLevel(ctx, UserPreloadMap, UserPayloadPreloadLevels.User)
	mods = append(mods, dm.UserWhere.ID.EQ(m.ID))
	pM, err := dm.Users(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.UserPayload{
		User: UserToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) CreateUsers(ctx context.Context, input fm.UsersCreateInput) (*fm.UsersPayload, error) {
	// TODO: Implement batch create
	return nil, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input fm.UserUpdateInput) (*fm.UserPayload, error) {
	m := UserUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)

	dbID := UserID(id)
	if _, err := dm.Users(
		dm.UserWhere.ID.EQ(dbID),
	).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	// resolve requested fields after updating
	mods := helper.GetPreloadModsWithLevel(ctx, UserPreloadMap, UserPayloadPreloadLevels.User)
	mods = append(mods, dm.UserWhere.ID.EQ(dbID))

	pM, err := dm.Users(mods...).One(context.Background(), boil.GetContextDB())
	return &fm.UserPayload{
		User: UserToGraphQL(pM, nil, 0),
	}, err
}

func (r *mutationResolver) UpdateUsers(ctx context.Context, filter *fm.UserFilter, input fm.UserUpdateInput) (*fm.UsersUpdatePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, UserFilterToMods(filter)...)

	m := UserUpdateInputToModelM(helper.GetInputFromContext(ctx, inputKey), input)
	if _, err := dm.Users(mods...).UpdateAll(context.Background(), boil.GetContextDB(), m); err != nil {
		return nil, err
	}

	return &fm.UsersUpdatePayload{
		Ok: true,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*fm.UserDeletePayload, error) {
	mods := []qm.QueryMod{
		dm.UserWhere.ID.EQ(UserID(id)),
	}
	_, err := dm.Users(mods...).DeleteAll(context.Background(), boil.GetContextDB())
	return &fm.UserDeletePayload{
		ID: id,
	}, err
}

func (r *mutationResolver) DeleteUsers(ctx context.Context, filter *fm.UserFilter) (*fm.UsersDeletePayload, error) {
	var mods []qm.QueryMod
	mods = append(mods, UserFilterToMods(filter)...)
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
	dbID := CommentID(id)
	mods := helper.GetPreloadMods(ctx, CommentPreloadMap)
	mods = append(mods, dm.CommentWhere.ID.EQ(dbID))
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	m, err := dm.Comments(mods...).One(context.Background(), boil.GetContextDB())
	return CommentToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Comments(ctx context.Context, filter *fm.CommentFilter, pagination *fm.CommentPagination) ([]*fm.Comment, error) {
	mods := helper.GetPreloadMods(ctx, CommentPreloadMap)
	mods = append(mods, dm.CommentWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, CommentPaginationToMods(pagination)...)
	mods = append(mods, CommentFilterToMods(filter)...)
	a, err := dm.Comments(mods...).All(context.Background(), boil.GetContextDB())
	return CommentsToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Company(ctx context.Context, id string) (*fm.Company, error) {
	dbID := CompanyID(id)
	mods := helper.GetPreloadMods(ctx, CompanyPreloadMap)
	mods = append(mods, dm.CompanyWhere.ID.EQ(dbID))
	m, err := dm.Companies(mods...).One(context.Background(), boil.GetContextDB())
	return CompanyToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Companies(ctx context.Context, filter *fm.CompanyFilter, pagination *fm.CompanyPagination) ([]*fm.Company, error) {
	mods := helper.GetPreloadMods(ctx, CompanyPreloadMap)
	mods = append(mods, CompanyPaginationToMods(pagination)...)
	mods = append(mods, CompanyFilterToMods(filter)...)
	a, err := dm.Companies(mods...).All(context.Background(), boil.GetContextDB())
	return CompaniesToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Follower(ctx context.Context, id string) (*fm.Follower, error) {
	dbID := FollowerID(id)
	mods := helper.GetPreloadMods(ctx, FollowerPreloadMap)
	mods = append(mods, dm.FollowerWhere.ID.EQ(dbID))
	m, err := dm.Followers(mods...).One(context.Background(), boil.GetContextDB())
	return FollowerToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Followers(ctx context.Context, filter *fm.FollowerFilter, pagination *fm.FollowerPagination) ([]*fm.Follower, error) {
	mods := helper.GetPreloadMods(ctx, FollowerPreloadMap)
	mods = append(mods, FollowerPaginationToMods(pagination)...)
	mods = append(mods, FollowerFilterToMods(filter)...)
	a, err := dm.Followers(mods...).All(context.Background(), boil.GetContextDB())
	return FollowersToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Location(ctx context.Context, id string) (*fm.Location, error) {
	dbID := LocationID(id)
	mods := helper.GetPreloadMods(ctx, LocationPreloadMap)
	mods = append(mods, dm.LocationWhere.ID.EQ(dbID))
	m, err := dm.Locations(mods...).One(context.Background(), boil.GetContextDB())
	return LocationToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Locations(ctx context.Context, filter *fm.LocationFilter, pagination *fm.LocationPagination) ([]*fm.Location, error) {
	mods := helper.GetPreloadMods(ctx, LocationPreloadMap)
	mods = append(mods, LocationPaginationToMods(pagination)...)
	mods = append(mods, LocationFilterToMods(filter)...)
	a, err := dm.Locations(mods...).All(context.Background(), boil.GetContextDB())
	return LocationsToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Post(ctx context.Context, id string) (*fm.Post, error) {
	dbID := PostID(id)
	mods := helper.GetPreloadMods(ctx, PostPreloadMap)
	mods = append(mods, dm.PostWhere.ID.EQ(dbID))
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	m, err := dm.Posts(mods...).One(context.Background(), boil.GetContextDB())
	return PostToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Posts(ctx context.Context, filter *fm.PostFilter, pagination *fm.PostPagination) ([]*fm.Post, error) {
	mods := helper.GetPreloadMods(ctx, PostPreloadMap)
	mods = append(mods, dm.PostWhere.UserID.EQ(
		auth.UserIDFromContext(ctx),
	))
	mods = append(mods, PostPaginationToMods(pagination)...)
	mods = append(mods, PostFilterToMods(filter)...)
	a, err := dm.Posts(mods...).All(context.Background(), boil.GetContextDB())
	return PostsToGraphQL(a, nil, 0), err
}

func (r *queryResolver) Role(ctx context.Context, id string) (*fm.Role, error) {
	dbID := RoleID(id)
	mods := helper.GetPreloadMods(ctx, RolePreloadMap)
	mods = append(mods, dm.RoleWhere.ID.EQ(dbID))
	m, err := dm.Roles(mods...).One(context.Background(), boil.GetContextDB())
	return RoleToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Roles(ctx context.Context, filter *fm.RoleFilter, pagination *fm.RolePagination) ([]*fm.Role, error) {
	mods := helper.GetPreloadMods(ctx, RolePreloadMap)
	mods = append(mods, RolePaginationToMods(pagination)...)
	mods = append(mods, RoleFilterToMods(filter)...)
	a, err := dm.Roles(mods...).All(context.Background(), boil.GetContextDB())
	return RolesToGraphQL(a, nil, 0), err
}

func (r *queryResolver) User(ctx context.Context, id string) (*fm.User, error) {
	dbID := UserID(id)
	mods := helper.GetPreloadMods(ctx, UserPreloadMap)
	mods = append(mods, dm.UserWhere.ID.EQ(dbID))
	m, err := dm.Users(mods...).One(context.Background(), boil.GetContextDB())
	return UserToGraphQL(m, nil, 0), err
}

func (r *queryResolver) Users(ctx context.Context, filter *fm.UserFilter, pagination *fm.UserPagination) ([]*fm.User, error) {
	mods := helper.GetPreloadMods(ctx, UserPreloadMap)
	mods = append(mods, UserPaginationToMods(pagination)...)
	mods = append(mods, UserFilterToMods(filter)...)
	a, err := dm.Users(mods...).All(context.Background(), boil.GetContextDB())
	return UsersToGraphQL(a, nil, 0), err
}

func (r *Resolver) Mutation() fm.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() fm.QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

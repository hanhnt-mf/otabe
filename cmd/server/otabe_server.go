package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	pbl_otabe "otabe"
	"otabe/controller"
	pb "otabe/pb"
	service "otabe/service"
	"time"
)


const (
	secretKey = "secret"
	tokenDuration = 15 * time.Minute
)
var (
	port = flag.Int("port", 8000, "The server port")
	db *sql.DB
	err error
	jwtManager = service.NewJWTManager(secretKey, tokenDuration)
)

type OTabeServer struct {
	pb.UnimplementedOTabeManagerServer
}

func GetMenuItems (menuId int32) ([]*pb.MenuItems, []*pb.NationsRate) {
	menuItems := make([]*pb.MenuItems, 0)
	userIds := make([]int32, 0)
	itemQuery, er := db.Query("SELECT id, name, description, price FROM item WHERE menu_id = ?", menuId)
	if er != nil {
		log.Fatalf("Cannot query to item with menuId reference %v", err)
		panic(er)
	}

	item := &pb.MenuItems{}
	itemIds := make([]int, 0)
	itemsNationRates := make([]*pb.NationsRate, 0)
	for itemQuery.Next() {
		itemId := new(int)
			err = itemQuery.Scan(&itemId, &item.ItemName, &item.Description, &item.Price)
			if err != nil {
				log.Fatalf("Cannot scan item query %v", err)
				panic(err)
			}
			itemIds = append(itemIds, *itemId)

			userItemFeedbackQuery, erF := db.Query("SELECT u.id, u.user_name, u.nation , if2.comment, if2.rate FROM user u JOIN item_feedback if2 ON u.id = if2.user_id  WHERE if2.item_id = ?", itemId)
			if erF != nil {
				log.Fatalf("Cannot query to feedback with itemId reference %v", err)
				panic(erF)
			}
			userComment := &pb.Comments{}
			itemRates := int32(0)
			itemNations := make([]string, 0)
			feedBacks := make([]*pb.Feedbacks, 0)
			usersInfo := &pb.Users{}
			for userItemFeedbackQuery.Next() {
				err = userItemFeedbackQuery.Scan(&userComment.UserId, &usersInfo.UserName, &usersInfo.Nation,	&userComment.UserComment, &userComment.Rate)
				if err != nil {
					log.Fatalf("Cannot scan item query %v", err)
					panic(err)
				}
				if !pbl_otabe.ItemExists(userIds, userComment.UserId) {
					userIds = append(userIds, userComment.UserId)
				}
				if pbl_otabe.ItemExists(itemNations, usersInfo.Nation) {
					itemRates = (itemRates + userComment.Rate)/2
					for i := range feedBacks {
						if feedBacks[i].Nation == usersInfo.Nation {
							feedBacks[i].Rate = itemRates
							feedBacks[i].Comments = append(feedBacks[i].Comments, &pb.Comments{UserId: userComment.UserId,
								UserComment: userComment.UserComment, Rate: userComment.Rate})
						}
					}
					for i := range itemsNationRates {
						if itemsNationRates[i].NationName == usersInfo.Nation {
							itemsNationRates[i].Rate = itemRates
						}
					}
				} else {
					itemNations = append(itemNations, usersInfo.Nation)
					itemRates = userComment.Rate
					userNewComments := make([]*pb.Comments, 0)
					userNewComments =  append(userNewComments, &pb.Comments{UserId: userComment.UserId, UserComment: userComment.UserComment,
						Rate: userComment.Rate})
					feedBacks = append(feedBacks, &pb.Feedbacks{Nation: usersInfo.Nation, Rate: itemRates, Comments: userNewComments})
					itemsNationRates = append(itemsNationRates, &pb.NationsRate{NationName: usersInfo.Nation, Rate: itemRates})
				}
			}
			menuItems = append(menuItems, &pb.MenuItems{ItemName: item.ItemName, Description: item.Description, Price: item.Price,
				Feedbacks: feedBacks})
	}

	return menuItems, itemsNationRates
}

func GetMenus(resId int32)  ([]*pb.Menus, []*pb.NationsRate) {
	// query to menus = resId -> item = menuId -> item_feedback = item_id -> user = user_id in item_feedback
	// -> get user_id: nation, rate -> push []userFeedbacks
	menuQuery, er := db.Query("SELECT id,name FROM menu WHERE restaurant_id = ?", resId)
	if er != nil {
		log.Fatalf("Cannot query to menu with resId reference %v", err)
		panic(er)
	}

	menu := &pb.Menus{}
	menus := make([]*pb.Menus, 0)
	nationsRate := make([]*pb.NationsRate, 0)
	for menuQuery.Next() {
		menuId := new(int32)
		err = menuQuery.Scan(menuId, &menu.Name)
		if err != nil {
			log.Fatalf("Cannot scan menu query %v", err)
			panic(err)
		}
		menuItems, menuRates := GetMenuItems(*menuId)
		menus = append(menus, &pb.Menus{Name: menu.Name, MenuItems: menuItems})
		if len(nationsRate) == 0 {
			nationsRate = menuRates
		} else {
			for menuIn := range menuRates {
				for nationIn := range nationsRate {
					if menuRates[menuIn].NationName == nationsRate[nationIn].NationName {
						nationsRate[nationIn].Rate = (nationsRate[nationIn].Rate + menuRates[menuIn].Rate)/2
					}
				}
			}
		}
	}
	return menus, nationsRate
}

func SearchRestaurants(req *pb.ListRestaurantsRequest) ([]int, error) {
	searchResSQL := `SELECT r.id, r.long, r.lat FROM restaurant as r 
    	INNER JOIN menu as m ON (r.id = m.restaurant_id AND (? IS NULL OR r.name = ?) AND (? IS NULL OR r.address = ?)) 
    	INNER JOIN item as i  ON (m.id = i.menu_id AND (? is NULL  OR i.name = ?) )`

	if req.Nation != nil {
		searchResSQL = fmt.Sprintf("%s INNER JOIN item_feedback as if2 ON (i.id = if2.item_id) ", searchResSQL)
		searchResSQL = fmt.Sprintf(`%s INNER JOIN user as u ON (if2.user_id = u.id AND ("%s" is NULL OR u.nation = "%s"))`, searchResSQL, *req.Nation, *req.Nation)
	}
	searchResSQL = fmt.Sprintf("%s  GROUP BY r.id ORDER BY r.%s DESC", searchResSQL, *req.SortedBy)
	searchResSQL = fmt.Sprintf("%s LIMIT %d OFFSET %d", searchResSQL, req.Paging.PageLimit, req.Paging.PageNumber - 1)

	searchResPrepare, errR := db.Prepare(searchResSQL)
	if errR != nil {
		log.Fatalf("Error preparing search restaurant sql %v", errR)
		panic(errR)
	}

	searchResQuery, errR := searchResPrepare.Query(
		req.RestaurantName, req.RestaurantName,
		req.Prefecture, req.Prefecture,
		req.ItemName, req.ItemName)

	if errR != nil {
		log.Fatalf("Error query search restaurant values %v", errR)
		panic(errR)
	}
	long := new(float64)
	lat := new(float64)
	restaurantIds := make([]int, 0)
	for searchResQuery.Next() {
		restaurantId := new(int)
		errR = searchResQuery.Scan(&restaurantId, &long, &lat)
		if errR != nil {
			log.Fatalf("Cannot scan restaurants list %v", err)
			panic(errR)
		}
		if req.Location != nil {
			pointsDistance := pbl_otabe.Distance(*lat, *long, *req.Location.Lat, *req.Location.Long)
			if pointsDistance <= 100000 {
				restaurantIds = append(restaurantIds, *restaurantId)
			}
		} else {
			restaurantIds = append(restaurantIds, *restaurantId)
		}
	}

	return restaurantIds, errR
}

func ConvertRestaurantConditions(req *pb.ListRestaurantsRequest) *pb.ListRestaurantsRequest {
	if req.GetPaging() == nil {
		req.Paging = &pb.Paging{PageLimit: uint64(10), PageNumber: uint64(1)}
	}
	if req.SortedBy == nil {
		defaultSortedBy := "created_at"
		req.SortedBy = &defaultSortedBy
	}
	if req.GetRestaurantName() == "" {
		req.RestaurantName = nil
	}
	if req.GetNation() == "" {
		req.Nation = nil
	}
	if req.GetPrefecture() == "" {
		req.Prefecture = nil
	}
	if req.GetItemName() == "" {
		req.ItemName = nil
	}
	return req
}

func (s *OTabeServer) ListRestaurantsByOptions(ctx context.Context, req *pb.ListRestaurantsRequest, ) (*pb.ListRestaurantsResponse, error) {
	errV := controller.ValidateListRestaurantsRequest(req)
	if errV != nil {
		return nil, errV
	}
	// convert condition : order by values
	convertedReqConditions := ConvertRestaurantConditions(req)

	restaurantIds, err := SearchRestaurants(convertedReqConditions)
	if err != nil {
		log.Fatalf("Err query searh restaurants by options %v", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	restaurantsList := make([]*pb.GetRestaurantResponse, 0)
	for _, restaurantId := range restaurantIds {
		restaurantDetails, err := s.GetRestaurantDetails(ctx, &pb.GetRestaurantRequest{RestaurantId: int32(restaurantId)})
		if err != nil {
			log.Fatalf("Err get restaurant details %v", err)
			panic(err)
		}
		restaurantsList = append(restaurantsList, restaurantDetails)
	}

	return &pb.ListRestaurantsResponse{Data: restaurantsList}, nil
}

// GetRestaurantDetails : api get restaurant details
func (s *OTabeServer) GetRestaurantDetails(ctx context.Context, req *pb.GetRestaurantRequest) (*pb.GetRestaurantResponse, error) {
	var restaurantId = req.GetRestaurantId()
	resQuery, er := db.Query("SELECT * FROM restaurant WHERE id = ?", restaurantId)
	if er != nil {
		log.Fatalf("Cannot get restaurant details %v", err)
		panic(er)
	}
	res := &pb.Restaurant{Geo: &pb.Geo{}}
	createdAt := make([]uint8, 0)
	updatedAt := make([]uint8, 0)
	for resQuery.Next() {
		err = resQuery.Scan(&res.Id, &res.Name, &res.Website, &res.Phone, &res.Description,
			&res.Address, &res.PostalCode, &res.Geo.Long, &res.Geo.Lat, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("Cannot get restaurant details %v", err)
			panic(err)
		}
	}


	menus, nationsRate := GetMenus(restaurantId)
	return &pb.GetRestaurantResponse{Restaurant: res, NationsRate: nationsRate, Menus: menus}, nil
}

func insertNewRestaurant(restaurant *pb.RestaurantRequest) (int32, error) {
	insertSQL := "INSERT INTO restaurant (name, website, phone, description, postal_code, address, `long`, lat, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	if restaurant.Id != nil {
		insertSQL = fmt.Sprintf("WHERE id = %d", *restaurant.Id)
	}

	insertResPrepare, errR := db.Prepare(insertSQL)
	if errR != nil {
		log.Fatalf("Error preparing insert restaurant sql %v", errR)
		panic(errR)
	}

	_, err := insertResPrepare.Query(
		restaurant.Name, restaurant.Website,restaurant.Phone, restaurant.Description,
		restaurant.PostalCode, restaurant.Address, restaurant.Geo.Long, restaurant.Geo.Lat,
		time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatalf("Err query create new restaurant %v", err)
		panic(err)
	}

	newRestaurant := &pb.Restaurant{}
	//for newRestaurantQuery.Next() {
	//	err = newRestaurantQuery.Scan(&newRestaurant.Id)
	//	if err != nil {
	//		log.Fatalf("Cannot scan restaurant details %v", err)
	//		panic(err)
	//	}
	//}
	restaurantQuery, err := db.Query("SELECT MAX(id) from restaurant")
	for restaurantQuery.Next() {
		err = restaurantQuery.Scan(&newRestaurant.Id)
		if err != nil {
			log.Fatalf("Cannot scan restaurant details %v", err)
			panic(err)
		}
	}

	restaurantId := newRestaurant.Id
	if restaurant.Id != nil {
		restaurantId = *restaurant.Id
	}

	return restaurantId, nil
}

func insertMenuItems(menus []*pb.MenuRequest, restaurantId int32) error {
	for _, menu := range menus {
		insertMenuSQL := `INSERT INTO menu (restaurant_id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
		if menu.Id != nil {
			insertMenuSQL = fmt.Sprintf("WHERE id = %d", *menu.Id)
		}
		insertMenuSQLPrepare, errR := db.Prepare(insertMenuSQL)
		if errR != nil {
			log.Fatalf("Error preparing insert menu sql %v", errR)
			panic(errR)
		}
		_, err := insertMenuSQLPrepare.Query(restaurantId, menu.Name,
			time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Fatalf("Err query create new menu %v", err)
			panic(err)
		}

		menuId := new(int32)
		if menu.Id != nil {
			menuId = menu.Id
		} else {
			menuQuery, err := db.Query("SELECT MAX(id) from menu")
			for menuQuery.Next() && menu.Id == nil {
				err = menuQuery.Scan(&menuId)
				if err != nil {
					log.Fatalf("Cannot scan menu id %v", err)
					panic(err)
				}
			}
		}
		log.Printf("menu %v", *menuId)

		for _, item := range menu.MenuItems {
			_, err := db.Query("INSERT INTO item (menu_id, name, description, price, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
				*menuId, item.ItemName, item.Description, item.Price, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Fatalf("Err query create new item %v", err)
				panic(err)
			}
		}
	}

	return nil
}
func (s *OTabeServer) CreateNewRestaurant(ctx context.Context, req *pb.CreateRestaurantRequest) (*pb.CreateRestaurantResponse, error) {
	// check geo
	existedRestaurant, err := db.Query("SELECT id FROM restaurant WHERE `long` = ? AND lat = ?",
		req.Restaurant.Geo.Long, req.Restaurant.Geo.Lat)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find restaurant: %v", err)

	}
	for existedRestaurant.Next() {
		return nil, status.Errorf(codes.InvalidArgument, "restaurant existed: %v", err)
	}

	restaurantId, err := insertNewRestaurant(req.Restaurant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error inserting create new restaurant : %v", err)
	}

	err = insertMenuItems(req.Menus, restaurantId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error inserting create new menu : %v", err)
	}

	resQuery, er := db.Query("SELECT * FROM restaurant WHERE id = ?", restaurantId)
	if er != nil {
		log.Fatalf("Cannot get restaurant details %v", err)
		panic(er)
	}
	res := &pb.Restaurant{Geo: &pb.Geo{}}
	createdAt := make([]uint8, 0)
	updatedAt := make([]uint8, 0)
	for resQuery.Next() {
		err = resQuery.Scan(&res.Id, &res.Name, &res.Website, &res.Phone, &res.Description,
			&res.Address, &res.PostalCode, &res.Geo.Long, &res.Geo.Lat, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("Cannot get restaurant details %v", err)
			panic(err)
		}
	}
	return &pb.CreateRestaurantResponse{Restaurant: res}, nil
}

func (s *OTabeServer) UpdateRestaurant(ctx context.Context, req *pb.CreateRestaurantRequest) (*pb.GetRestaurantResponse, error) {
	//  find res -> update info res
	// find menus -> update each
	// find items -> update each
	var restaurantId = req.Restaurant.GetId()
	resQuery, er := db.Query("SELECT * FROM restaurant WHERE id = ?", restaurantId)
	if er != nil {
		log.Fatalf("Cannot get restaurant details %v", err)
		panic(er)
	}
	if !resQuery.Next() {
		return nil, controller.RestaurantNotFound()
	}

	restaurantId, err := insertNewRestaurant(req.Restaurant)
	if err != nil {
		log.Fatalf("Error inserting updating restaurant %v", err)
		panic(err)
	}

	err = insertMenuItems(req.Menus, restaurantId)
	if err != nil {
		log.Fatalf("Error inserting updating menus items %v", err)
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	restaurantDetails, err := s.GetRestaurantDetails(ctx, &pb.GetRestaurantRequest{RestaurantId: restaurantId})
	if err != nil {
		log.Fatalf("Error getting res details %v", err)
		panic(err)
	}
	return restaurantDetails, nil
}

func (s *OTabeServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// find user with username
	userQuery, err := db.Query("SELECT user_name, password, role FROM user WHERE user_name = ?", req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	user := &service.User{}
	for userQuery.Next() {
		err = userQuery.Scan(&user.Username, &user.HashedPassword, &user.Role)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot scan user: %v", err)
		}
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.LoginResponse{AccessToken: token}
	return res, nil
}


func connect() {
	db, err = sql.Open("mysql", "docker:Hannamysql.1518@tcp(db:3306)/otabe")
	if err != nil {
		log.Fatalf("Error validating sql.Open arguments")
		panic(err)
	}

	//err = db.Ping()
	//if err != nil {
	//	log.Fatalf("Error verifying with db.Ping %v", err)
	//	panic(err)
	//}
}

func unaryServerInterceptor (
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func accessibleRoles() map[string][]string {
	const restaurantsPath = "/v1.OTabeManager/"
	return map[string][]string{
		restaurantsPath + "CreateNewRestaurant" : {"admin"},
		restaurantsPath + "UpdateNewRestaurant" : {"admin"},
		//restaurantsPath + "GetRestaurantDetails" : {"admin", "user"},
		//restaurantsPath + "ListRestaurantsByOptions" : {"admin", "user"},
	}
}

func main() {
	lis, er := net.Listen("tcp",":8080")
	if er != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()))
	pb.RegisterOTabeManagerServer(grpcServer, &OTabeServer{})
	log.Printf("Server listening at port %v", lis.Addr())

	// connect to database
	connect()

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
	defer db.Close()
}
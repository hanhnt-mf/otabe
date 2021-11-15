package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
	pbl_otabe "otabe"
	pb "otabe/v1"
)

var (
	port = flag.Int("port", 8000, "The server port")
	db *sql.DB
	err error
	nationsRate []*pb.NationsRate
)

type oTabeServer struct {
	pb.UnimplementedOTabeManagerServer
}

func getMenuItems (menuId int32) ([]*pb.MenuItems, []*pb.NationsRate) {
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

			userItemFeedbackQuery, erF := db.Query("SELECT u.id, u.first_name, u.last_name, u.nation , if2.comment, if2.rate FROM user u JOIN item_feedback if2 ON u.id = if2.user_id  WHERE if2.item_id = ?", itemId)
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
				err = userItemFeedbackQuery.Scan(&userComment.UserId, &usersInfo.FirstName,
					&usersInfo.LastName, &usersInfo.Nation,	&userComment.UserComment, &userComment.Rate)
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

	log.Printf("item nations rates %v", itemsNationRates)
	return menuItems, itemsNationRates
}

func getMenus(resId int32)  []*pb.Menus {
	// query to menus = resId -> item = menuId -> item_feedback = item_id -> user = user_id in item_feedback
	// -> get user_id: nation, rate -> push []userFeedbacks
	menuQuery, er := db.Query("SELECT id,name FROM menu WHERE restaurant_id = ?", resId)
	if er != nil {
		log.Fatalf("Cannot query to menu with resId reference %v", err)
		panic(er)
	}

	menu := &pb.Menus{}
	menus := make([]*pb.Menus, 0)
	for menuQuery.Next() {
		menuId := new(int32)
		err = menuQuery.Scan(menuId, &menu.Name)
		if err != nil {
			log.Fatalf("Cannot scan menu query %v", err)
			panic(err)
		}
		menuItems, menuRates := getMenuItems(*menuId)
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
	return menus
}

func (s *oTabeServer) GetRestaurantDetails(ctx context.Context, req *pb.GetRestaurantRequest) (*pb.GetRestaurantResponse, error) {
	log.Printf("Received: %v", req.GetRestaurantId())
	var restaurantId = req.GetRestaurantId()
	resQuery, er := db.Query("SELECT * FROM restaurant WHERE id = ?", restaurantId)
	if er != nil {
		log.Fatalf("Cannot get restaurant details %v", err)
		panic(er)
	}
	res := &pb.Restaurant{Geo: &pb.Geo{}}
	for resQuery.Next() {
		err = resQuery.Scan(&res.Id, &res.Name, &res.Website, &res.Phone, &res.Description,
			&res.Address, &res.PostalCode, &res.Geo.Long, &res.Geo.Lat, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			log.Fatalf("Cannot get restaurant details %v", err)
			panic(err)
		}
	}
	return &pb.GetRestaurantResponse{Restaurant: res, NationsRate: nationsRate, Menus: getMenus(restaurantId)}, nil
}


func connect() {
	db, err = sql.Open("mysql", "root:Hannamysql.1518@tcp(127.0.0.1:50125)/otabe")
	if err != nil {
		log.Fatalf("Error validating sql.Open arguments")
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error verifying with db.Ping %v", err)
		panic(err)
	}
}
func main() {
	lis, er := net.Listen("tcp",":8080")
	if er != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOTabeManagerServer(grpcServer, &oTabeServer{})
	log.Printf("Server listening at port %v", lis.Addr())

	// connect to database
	connect()
	//getMenus(1)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
	defer db.Close()

}
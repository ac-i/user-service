package store

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ac-i/user-service/config"
	"github.com/ac-i/user-service/proto/model"
	"github.com/ac-i/user-service/proto/serv"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserStore struct {
	// dev|prod
	mode string
	// Memory
	// dbM *StoreDB
	// User
	dbU *StoreDB
}

var UnimplementedUserServer serv.UnimplementedUserServer

func (store *UserStore) Open() *UserStore {
	modeF := flag.String("mode", "dev", "runtime mode: dev, prod")
	flag.Parse()

	// store := new(UserStore)

	// configuration
	store.mode = *modeF

	// M U separate dedicated storage space (M in memory)
	// store.dbM = store.initDB("M")
	store.dbU = store.initDB("U")

	// initial maintenance
	store.checkDB(true, true)

	return store
}

func (us *UserStore) initDB(key string) *StoreDB {
	// log.Println("initStore :: ", key, s.mode)

	var (
		tmp, mp string
		err     error
	)
	store := new(StoreDB)

	store.cfg.key = key
	store.cfg.mode = us.mode

	if store.cfg.mode == "dev" {
		mp = "d"
	} else if store.cfg.mode == "prod" {
		mp = "p"
	} else {
		err = errors.New("ERR :: initStore :: unknown mode: " + store.cfg.mode)
		log.Fatal(err)
	}

	// U X
	if config.ServDev.PresistantStore {
		// ! config.ServDev.PresistantStore
		tmp, err = os.UserCacheDir() //+
		if err != nil {
			log.Fatal("tmp", tmp, err)
		}
		tmp = tmp + string(os.PathSeparator) + "user-service" + string(os.PathSeparator)
		err = os.MkdirAll(tmp, 0766)
		if err != nil {
			log.Fatal(err)
		}
	}

	if store.cfg.key == "M" {
		store.cfg.dbpath = ":memory:"
		store.cfg.isPresistant = false
	} else if store.cfg.key == "U" {
		// ! config.ServDev.PresistantStore
		if config.ServDev.PresistantStore {
			store.cfg.dbpath = tmp + "store_u_" + mp + "_.db"
			store.cfg.isPresistant = true
		} else {
			store.cfg.dbpath = ":memory:"
			store.cfg.isPresistant = false
		}
		log.Println("User Store: ", store.cfg.dbpath)
	} else {
		err = errors.New("ERR :: initStore :: unknown key: " + store.cfg.key)
		log.Fatal(err)
	}

	store.db, err = buntdb.Open(store.cfg.dbpath)
	if err != nil {
		log.Fatal("ERR :: initStore :: buntdb.Open :: ", store.cfg.dbpath, err)
	}

	return store
}

func (us *UserStore) checkDB(doShrink, doIndex bool) {
	if doShrink {
		_ = us.dbU.ShrinkDB()
		// _ = us.dbM.ShrinkDB()
	}

	if doIndex {
		if err := us.dbU.IndexJSON("userid", "*", "userid"); err != nil {
			log.Println("ERR :: createIndexJSON :: ", err)
		}
		if err := us.dbU.IndexJSON("active", "*", "active"); err != nil {
			log.Println("ERR :: createIndexJSON :: ", err)
		}

		//s.ensureCollateJSON("qv", "*", "ENGLISH_LOOSE", "qv")
		//s.ensureCollateJSON("qx", "*", "ENGLISH_LOOSE", "qx")
		//s.db.CreateIndex("last_name_age", "*", buntdb.IndexJSON("name.last"), buntdb.Desc(buntdb.IndexJSON("age")))
	}
}

// DONE UserAdd
func (us *UserStore) UserAdd(ctx context.Context, in *model.User) (out *model.User, err error) {
	_, _ = UnimplementedUserServer.Add(ctx, in)
	if in.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument: user name is empty")
	}
	// check and calculate next UsedID
	// TODO use singular read/write transaction to calc next id and write at once for consistency at high loads
	in.Userid = us.nextUserID()

	// jb, err := json.Marshal(&in)
	jb, err := in.DoMarshal()
	if err != nil {
		return out, err
	}

	// store User in dbU, id prefix: u:
	kvU := us.dbU.SetKV(fmt.Sprintf("u:%d", in.Userid), string(jb), 0)

	err = json.Unmarshal([]byte(kvU.value), &out)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (us *UserStore) maxUserID() (x int32) {
	// return us.dbU.txMaxUserID()
	tx, err := us.dbU.db.Begin(false)
	if err != nil {
		return x
	}
	err = tx.Descend("userid", func(key, value string) bool {
		vx := int32(gjson.Get(value, "userid").Int())
		if vx > x {
			x = vx
		}
		return false
	})
	if err != nil {
		return x
	}
	if err := tx.Rollback(); err != nil {
		return x
	}
	return x
}

func (us *UserStore) nextUserID() (x int32) {
	// TODO deprecate: use singular read/write transaction in UserAdd
	x = us.maxUserID()
	if x <= 1000 {
		return 1000 + 1
	} else {
		return x + 1
	}

}

// DONE UserGetList
func (us *UserStore) UserGetList(ctx context.Context, in *model.UserSelect) (*model.Users, error) {
	_, _ = UnimplementedUserServer.GetList(ctx, in)
	// using string to simulate nil of a bool
	// active=false|true|'' - to filter false or true or all
	// PLUS DEV MODE TESTS: active=dev-db-mockup|dev-db-truncate - for tests in dev mode
	// a=0 false, a=1 true, a='' both, other reject
	var a int
	t := strings.ToLower(in.Active)
	if t == "" {
		a = 2
	} else if t == "0" || t == "false" {
		a = 0
	} else if t == "1" || t == "true" {
		a = 1
	} else if t == "dev-db-shrink" {
		_ = us.dbU.ShrinkDB()
		a = 2
	} else if t == "dev-db-mockup" {
		// _ = us.dbU.ShrinkDB()
		us.mockupUsers(10, false)
		a = 2
	} else if t == "dev-db-mockup-100" {
		// _ = us.dbU.ShrinkDB()
		us.mockupUsers(100, false)
		a = 2
	} else if t == "dev-db-truncate" {
		_ = us.dbU.TruncateDB()
		_ = us.dbU.ShrinkDB()
		a = 2
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "invalid query parameter")
	}
	return us.dbU.txUserGetList(a)
}

// DONE txUserGetList
func (st *StoreDB) txUserGetList(a int) (*model.Users, error) {
	// a=0 false, a=1 true, a='' both, other reject
	tx, err := st.db.Begin(false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	out := new(model.Users)
	err = tx.Ascend("userid", func(key, value string) bool {
		vb := []byte(value)
		// skip records witout userid - just in case
		if gjson.GetBytes(vb, "userid").Int() == 0 {
			return true
		}
		ac := gjson.GetBytes(vb, "active").Bool()
		var add bool
		if a == 2 {
			add = true
		} else if ac && (a == 1) {
			add = true
		} else if !ac && a == 0 {
			add = true
		} else {
			add = false
		}
		if add {
			user := new(model.User)
			err = json.Unmarshal(vb, &user)
			if err == nil {
				out.Users = append(out.Users, user)
			}
		}
		return true
	})
	if err != nil {
		log.Println(err)
	}
	if err := tx.Rollback(); err != nil {
		log.Println(err)
	}
	// return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
	return out, nil
}

// DONE Del - gw delete: "/users/{userid}"
func (us *UserStore) UserDel(ctx context.Context, in *model.User) (*model.User, error) {
	_, _ = UnimplementedUserServer.Del(ctx, in)

	if in.Userid <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument: userid")
	}

	out := new(model.User)
	kvU := us.dbU.DelKV(fmt.Sprintf("u:%d", in.Userid))
	if kvU.err != nil {
		// TODO confirm: no error if user id is already deleted
		// return nil, kvU.err
		return out, nil
	} //else {
	// _ = us.dbU.ShrinkDB()
	// }

	err := json.Unmarshal([]byte(kvU.old), &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (us *UserStore) mockupUsers(max int, randomActive bool) {
	if us.maxUserID() == 0 {
		rand.Seed(time.Now().UnixNano())
		var active bool
		var x int
		for _, f1 := range strings.Split("Peter Georgia Mark Evelyn Ella Madison Julian Lincoln Christian Lucy", " ") {
			// for _, f2 := range strings.Split("Samson Carol Alan Tom Melinda Janet Donald Samuel George John", " ") {
			for _, l1 := range strings.Split("Anderson Cooper Johnson Cohen Brown Tesla Wells Orwell Snow Cage", " ") {
				// for _, l2 := range strings.Split("Smith Davies Taylor Walker Robinson Wood Martin Clarke Green White", " ") {
				x += 1
				if x > max || x > 100 {
					return
				}
				user := new(model.User)
				if randomActive {
					user.Active = (rand.Intn(100) >= 50)
				} else {
					if active {
						active = false
					} else {
						active = true
					}
					user.Active = active
				}
				user.Name = f1 + " " + l1

				_, err := us.UserAdd(context.TODO(), user)
				if err != nil {
					log.Println(err)
					return
					// return nil, err
				}
				// }
			}
			// }
		}
	}
}

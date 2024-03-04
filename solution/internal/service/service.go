package service

import (
	"context"
	"fmt"
	"solution/internal/contract"
	"solution/internal/utils"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool  *pgxpool.Pool
	pgUrl string
}

func New(pgUrl string) (*Service, error) {
	s := &Service{pgUrl: pgUrl}
	err := s.connect()
	if err != nil {
		return s, err
	}
	err = s.initTables(
		CountryDriver{},
		UserDriver{},
		RelationDriver{},
		PostDriver{},
		ReactionDriver{},
	)
	return s, err
}

func (s *Service) connect() error {
	pool, err := pgxpool.New(context.Background(), s.pgUrl)
	if err != nil {
		return err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return err
	}
	s.pool = pool
	return nil

}

func (s *Service) initTables(creators ...TableCreater) error {
	for _, creator := range creators {
		err := InitTable(creator, s.pool)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetCountries() ([]contract.Country, error) {
	return QueryAll(CountryDriver{}, s.pool, `SELECT * FROM countries`)
}

func (s *Service) GetCountriesOfRegions(regions []string) ([]contract.Country, error) {
	regionsRaw := strings.Join(utils.Map(regions, func(r string) string { return `'` + r + `'` }), ", ")
	return QueryAll(CountryDriver{}, s.pool, fmt.Sprintf(`SELECT * FROM countries WHERE region in (%v)`, regionsRaw))
}

func (s *Service) GetCountryByAlpha2(alpha2 string) (contract.Country, error) {
	query := fmt.Sprintf(`SELECT * FROM countries WHERE alpha2='%v'`, alpha2)
	return QuerySingle(CountryDriver{}, s.pool, query)
}

func (s *Service) GetUserByLogin(login string) (contract.User, error) {
	query := fmt.Sprintf(`SELECT * FROM users WHERE login='%v'`, login)
	return QuerySingle(UserDriver{}, s.pool, query)
}

func (s *Service) UserExists(user contract.User) bool {
	query := fmt.Sprintf(`SELECT * FROM users WHERE login='%v' OR email='%v' OR (phone='%v' AND phone != '')`, user.Login, user.Email, user.Phone)
	users, _ := QueryAll(UserDriver{}, s.pool, query)
	return len(users) > 0
}

func (s *Service) UserDataExists(user contract.User) bool {
	query := fmt.Sprintf(`SELECT * FROM users WHERE ((phone='%v' AND phone != '') OR email='%v') AND login !='%v'`, user.Phone, user.Email, user.Login)
	users, _ := QueryAll(UserDriver{}, s.pool, query)
	return len(users) > 0
}

func (s *Service) AddUser(user contract.User) error {
	return Insert(UserDriver{}, s.pool, user)
}

func (s *Service) UpdateUser(old, newUser contract.User) error {
	return Update(UserDriver{}, s.pool, old, newUser)
}

func (s *Service) FindRelation(senderLogin, accepterLogin string) (contract.Relation, error) {
	query := fmt.Sprintf(`SELECT * FROM relations WHERE senderLogin='%v' AND accepterLogin='%v'`, senderLogin, accepterLogin)
	return QuerySingle(RelationDriver{}, s.pool, query)
}

func (s *Service) AddToFriends(senderLogin, accepterLogin string) error {
	return Insert(RelationDriver{}, s.pool, contract.Relation{
		SenderLogin:   senderLogin,
		AccepterLogin: accepterLogin,
		CreateTime:    time.Now().Unix(),
	})
}

func (s *Service) DeleteRelation(relation contract.Relation) error {
	return Delete(RelationDriver{}, s.pool, relation)
}

func (s *Service) GetFriends(login string, limit, offset int) ([]AccepterRelation, error) {
	query := fmt.Sprintf(`SELECT * FROM relations JOIN users on users.login=relations.accepterLogin WHERE senderLogin='%v' ORDER BY -createTime LIMIT %v OFFSET %v`, login, limit, offset)
	rows, err := QueryAll(AccepterRelationDriver{}, s.pool, query)
	if err != nil {
		return []AccepterRelation{}, err
	}
	return rows, nil
}

func (s *Service) GetUsers() ([]contract.User, error) {
	return QueryAll(UserDriver{}, s.pool, `SELECT * FROM users`)
}

func (s *Service) GetPosts() ([]contract.Post, error) {
	return QueryAll(PostDriver{}, s.pool, `SELECT * FROM posts`)
}

func (s *Service) AddPost(p contract.Post) error {
	return Insert(PostDriver{}, s.pool, p)
}

func (s *Service) GetPostById(id string) (contract.Post, error) {
	post, err := QuerySingle(PostDriver{}, s.pool, fmt.Sprintf(`SELECT * FROM posts WHERE id='%v'`, id))
	if err != nil {
		return contract.Post{}, err
	}
	likes, dislikes, err := s.GetReactionsCount(id)
	if err != nil {
		return contract.Post{}, err
	}
	post.LikesCount = likes
	post.DislikesCount = dislikes
	return post, nil
}

func (s *Service) GetReactionsCount(postId string) (int, int, error) {
	query := fmt.Sprintf(`SELECT * FROM reactions WHERE postId='%v'`, postId)
	reactions, err := QueryAll(ReactionDriver{}, s.pool, query)
	if err != nil {
		return 0, 0, err
	}
	return len(utils.Filter(reactions, func(r contract.Reaction) bool {
			return r.Type == contract.Like
		})), len(utils.Filter(reactions, func(r contract.Reaction) bool {
			return r.Type == contract.Dislike
		})), nil
}

func (s *Service) GetPostsOf(author string, limit, offset int) ([]contract.Post, error) {
	posts, err := QueryAll(PostDriver{}, s.pool, fmt.Sprintf(
		`SELECT * FROM posts WHERE author='%v'
		ORDER BY -createdAt
		LIMIT %v OFFSET %v`,
		author, limit, offset))
	if err != nil {
		return []contract.Post{}, err
	}
	posts = utils.Map(posts, func(p contract.Post) contract.Post {
		likes, dislikes, e := s.GetReactionsCount(p.Id)
		fmt.Println(p.Id, likes, dislikes, e)
		if e != nil {
			err = e
		}
		p.LikesCount = likes
		p.DislikesCount = dislikes
		return p
	})
	return posts, err
}

func (s *Service) FindReaction(login, postId string) (contract.Reaction, error) {
	return QuerySingle(ReactionDriver{}, s.pool, fmt.Sprintf(`SELECT * FROM reactions WHERE login='%v' AND postId='%v'`,
		login, postId))
}

func (s *Service) SetReaction(login, postId string, rType contract.ReactionType) error {
	react, err := s.FindReaction(login, postId)
	if err == nil {
		newReact := react
		newReact.Type = rType
		return Update(ReactionDriver{}, s.pool, react, newReact)
	}
	return Insert(ReactionDriver{}, s.pool, contract.Reaction{
		Login:  login,
		PostId: postId,
		Type:   rType,
	})
}

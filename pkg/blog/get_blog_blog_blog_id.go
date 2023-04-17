package blog

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetBlogBlogBlogId(ctx echo.Context, dbClient db.Client, blogId string) (api.ResGetBlogBlogBlogId, *response.Error) {
	res := api.ResGetBlogBlogBlogId{}
	blog, err := getBlogFromTagId(dbClient, blogId)
	if err != nil {
		return api.ResGetBlogBlogBlogId{}, err
	}
	rerr := copier.Copy(&res, &blog)
	if rerr != nil {
		return api.ResGetBlogBlogBlogId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました4", Log: rerr.Error()}
	}
	return res, nil
}

type blog struct {
	Author    blogObjectAuthor
	Name      string `db:"name"`
	Tags      []blogObjectTag
	BlogId    string `db:"blog_id"`
	Content   string `db:"content"`
	IsPublic  string `db:"is_public"`
	TopImage  string `db:"top_image"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type blogObjectAuthor struct {
	UserId   string `db:"user_id"`
	Username string `db:"username"`
	IconUrl  string `db:"icon_url"`
}

type blogObjectTag struct {
	TagId string `db:"tag_id"`
	Name  string `db:"name"`
}

func getBlogFromTagId(dbClient db.Client, blogId string) (blog, *response.Error) {
	params := struct {
		BlogId string `twowaysql:"blogId"`
	}{
		BlogId: blogId,
	}
	blogs := []blog{}
	rerr := dbClient.Select(&blogs, "sql/blog/select_blog_from_blog_id.sql", &params)
	fmt.Printf("%+v", params)
	if rerr != nil {
		return blog{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました3", Log: rerr.Error()}
	}
	if len(blogs) == 0 {
		return blog{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "作品が存在しません", Log: "no rows in result"}
	}
	blogAuthor, err := getBlogBlogAuthorId(dbClient, blogId)
	if err != nil {
		return blog{}, err
	}
	blogs[0].Author = blogAuthor
	blogTags, err := getBlogBlogTagList(dbClient, blogId)
	if err != nil {
		return blog{}, err
	}
	blogs[0].Tags = blogTags
	return blogs[0], nil
}

func getBlogBlogTagList(dbClient db.Client, blogId string) ([]blogObjectTag, *response.Error) {
	params := struct {
		BlogId string `twowaysql:"blogId"`
	}{
		BlogId: blogId,
	}
	blogTags := []blogObjectTag{}
	err := dbClient.Select(&blogTags, "sql/blog/select_blog_blog_tag.sql", &params)
	if err != nil {
		return []blogObjectTag{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました1", Log: err.Error()}
	}
	return blogTags, nil
}

func getBlogBlogAuthorId(dbClient db.Client, blogId string) (blogObjectAuthor, *response.Error) {
	params := struct {
		BlogId string `twowaysql:"blogId"`
	}{
		BlogId: blogId,
	}
	blogAuthor := []blogObjectAuthor{}
	err := dbClient.Select(&blogAuthor, "sql/blog/select_blog_blog_author.sql", &params)
	if err != nil {
		return blogObjectAuthor{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました2", Log: err.Error()}
	}
	return blogAuthor[0], nil
}

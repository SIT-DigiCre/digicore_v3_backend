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

func GetBlogBlog(ctx echo.Context, dbClient db.Client, params api.GetBlogBlogParams) (api.ResGetBlogBlog, *response.Error) {
	res := api.ResGetBlogBlog{}
	blog, err := getBlogList(dbClient, params.Offset, params.AuthorId)
	if err != nil {
		return api.ResGetBlogBlog{}, err
	}
	rerr := copier.Copy(&res.Blogs, &blog)
	if rerr != nil {
		return api.ResGetBlogBlog{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type blogOverview struct {
	Author    blogObjectAuthor
	Name      string `db:"name"`
	Tags      []blogObjectTag
	BlogId    string `db:"blog_id"`
	Abstract  string `db:"abstract"`
	IsPublic  string `db:"is_public"`
	TopImage  string `db:"top_image"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func getBlogList(dbClient db.Client, offset *int, authorId *string) ([]blogOverview, *response.Error) {
	params := struct {
		Offset   *int    `twowaysql:"offset"`
		AuthorId *string `twowaysql:"authorId"`
	}{
		Offset:   offset,
		AuthorId: authorId,
	}
	blogOverviews := []blogOverview{}
	err := dbClient.Select(&blogOverviews, "sql/blog/select_blog.sql", &params)
	fmt.Printf("%+v", err)
	if err != nil {
		return []blogOverview{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました1", Log: err.Error()}
	}
	if len(blogOverviews) == 0 {
		return []blogOverview{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "作品が存在しません", Log: "no rows in result"}
	}
	for i := range blogOverviews {
		blogId := blogOverviews[i].BlogId
		blogAuthor, err := getBlogBlogAuthorId(dbClient, blogId)
		if err != nil {
			return []blogOverview{}, err
		}
		blogOverviews[i].Author = blogAuthor
		blogTags, err := getBlogBlogTagList(dbClient, blogId)
		if err != nil {
			return []blogOverview{}, err
		}
		blogOverviews[i].Tags = blogTags
	}
	return blogOverviews, nil
}

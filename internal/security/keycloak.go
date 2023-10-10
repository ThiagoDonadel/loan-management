package security

import (
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ThiagoDonadel/loan-management/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.ReplaceAll(ctx.Request.Header.Get("Authorization"), "Bearer ", "")

		config := viper.Sub("security.keycloak")
		realm := config.GetString("realm")

		client := gocloak.NewClient(config.GetString("hostname"))

		tokenValidation, err := client.RetrospectToken(ctx, token, config.GetString("client-id"), config.GetString("client-secret"), realm)

		if err != nil || !*tokenValidation.Active {
			ctx.JSON(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		_, claims, err := client.DecodeAccessToken(ctx, token, realm)

		if err != nil {
			ctx.JSON(http.StatusForbidden, "Forbidden")
			ctx.Abort()
			return
		}

		userToken := (*claims)["sub"].(string)
		userRequest := ctx.Param(utils.OWNER_ID_PARAM_NAME)

		if userToken != userRequest {
			roles := (*claims)["realm_access"].(map[string]any)["roles"].([]any)
			if !checkForRole("loan-admin", roles) {
				ctx.JSON(http.StatusForbidden, "Forbidden")
				ctx.Abort()
				return
			}

		}

	}
}

func checkForRole(targetRole string, roles []any) bool {

	for _, role := range roles {
		if role == targetRole {
			return true
		}
	}

	return false
}

/*

	fmt.Println(viper.Get("security.keycloak"))

	err := viper.UnmarshalKey("security.keycloak", &config)

	if err != nil {
		panic(err)
	}

	fmt.Println(config)


	ctx := context.Background()

	token := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJaMi10YzlrUUtEb3VZU2xqM25HZENza0dWOXo1c2pEaUVaM0VUWjJCTlpJIn0.eyJleHAiOjE2OTYzMzc4NTEsImlhdCI6MTY5NjMzNzU1MSwianRpIjoiYzkxNDFmMjItMGM0Zi00N2NhLTk1NmItMjJiMThjYmI1MTdjIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4OTg5L3JlYWxtcy9sb2FuIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6ImI3ZjE1ODE0LTI3NmEtNDM1Mi04ZTNhLTViOGM4MWY1MDY4NSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImE4ZDI1NDhkLWY1YTgtNGEwMi1iMTgxLTRlMWY0ZTc4ZDZjNSIsInNlc3Npb25fc3RhdGUiOiI2NjNkYWQ5Yy1iZDNkLTQyZDgtODY3Ny04YjY4YTg2MGMzYmIiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1sb2FuIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6ImVtYWlsIHByb2ZpbGUiLCJzaWQiOiI2NjNkYWQ5Yy1iZDNkLTQyZDgtODY3Ny04YjY4YTg2MGMzYmIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6IlRlc3RlIFRlc3RlIiwicHJlZmVycmVkX3VzZXJuYW1lIjoidGVzdCIsImdpdmVuX25hbWUiOiJUZXN0ZSIsImZhbWlseV9uYW1lIjoiVGVzdGUiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20ifQ.F-iR8H6lj7WmpANl697rkFUjsUuuBIllQ2xoMi0H_0mmZg5wYNe1VqkJBIqAo4NYjf2XuG3jrJllF1sJIS9Ede_4ZN0iJhfW9DT6jYyVTgnzLDcwyXcDRD09AakzWDzStrkcx9cwrjPU8haAAbRRdeb_A9diSPRK8Iumv6wYATk3sjlE91mlzn5pkZkhYhBV6mQbfk7EB7gs_Jl18omj1xg16kw3M_j8VpylHwRDB_0LeK9NvhVv6ItxfYTxgvby1ywRQCXm_EO8JtqN0NikVKl6gxACmzW_yPBDaGZZVnnUO8sXA4W2ISWoEBR9u2E60OYE14zjure7Eqvhv89EQw"

	rptResult, err := client.RetrospectToken(ctx, token, config.ClientId, config.ClientSecret, config.Realm)

	if err != nil {
		panic(err)
	}

	fmt.Println(rptResult)
} */

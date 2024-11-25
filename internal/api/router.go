package api

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(handler *Handler) *gin.Engine {
    router := gin.Default()

    // Swagger documentation
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // API routes
    v1 := router.Group("/api/v1")
    {
        songs := v1.Group("/songs")
        {
            songs.POST("", handler.CreateSong)
            songs.GET("", handler.ListSongs)
            songs.GET("/:id", handler.GetSong)
            songs.GET("/:id/verses", handler.GetSongVerses)
            songs.PUT("/:id", handler.UpdateSong)
            songs.DELETE("/:id", handler.DeleteSong)
        }
    }

    return router
}

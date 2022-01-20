package apimgr

import (
	"context"
	"reflect"
	"time"
	"unsafe"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	config "github.com/NpoolPlatform/go-service-framework/pkg/config"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	logger "github.com/NpoolPlatform/go-service-framework/pkg/logger"
	apimgr "github.com/NpoolPlatform/message/npool/apimgr"
)

func reliableRegister(apis *api.ServiceApis) {
	for {
		conn, err := grpc2.GetGRPCConn("api-manager.npool.top", grpc2.GRPCTAG)
		if err != nil {
			logger.Errorf("fail get api manager connection: %v", err)
			time.Sleep(time.Minute)
			continue
		}

		cli := goodspb.NewApiManagerClient(conn)

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

		err = cli.Register(ctx, &apimgr.RegisterRequest{
			Info: apis,
		})
		if err == nil {
			return
		}

		logger.Errorf("fail register apis: %v", err)
		time.Sleep()

		cancel()
		conn.Close()
	}
}

func Register(mux *runtime.ServeMux) {
	apis := &apimgr.ServiceApis{
		ServiceName: config.GetStringValueWithNameSpace("", config.KeyHostname),
	}

	valueOfMux := reflect.ValueOf(mux).Elem()
	handlers := valueOfMux.FieldByName("handlers")
	methIter := handlers.MapRange()
	for methIter.Next() {
		for i := 0; i < methIter.Value().Len(); i++ {
			pat := methIter.Value().Index(i).FieldByName("pat")
			tmp := reflect.NewAt(pat.Type(), unsafe.Pointer(pat.UnsafeAddr())).Elem()
			str := tmp.MethodByName("String").Call(nil)[0].String()
			apis.Paths = append(apis.Paths, &apimgr.Path{
				Method: methIter.Key().String(),
				Path:   str,
			})
		}
	}

	go reliableRegister(apis)
}

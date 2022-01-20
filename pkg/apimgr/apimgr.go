package apimgr

import (
	"reflect"
	"unsafe"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	config "github.com/NpoolPlatform/go-service-framework/pkg/config"
	apimgr "github.com/NpoolPlatform/message/npool/apimgr"
)

func Register(mux *runtime.ServeMux) {
	apis := apimgr.ServiceApis{
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
}

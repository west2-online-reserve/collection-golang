package config
import(
	"github.com/spf13/viper"
)
var Vp *viper.Viper
func InitConfig()(err error){
	Vp =viper.New()
	Vp.SetConfigFile("config/config.yaml")
	err = Vp.ReadInConfig()
	if err!= nil{
		return err
	}	
	return nil
}
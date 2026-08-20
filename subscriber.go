package main
import(
	"fmt"
)
func startSubscriber(){
	sub:=rdb.Subscribe(ctx, "link:clicks")
	ch:=sub.Channel()

	fmt.Println("subscriber started, lsitening for clicks")

	for msg:=range ch{
		fmt.Printf("link clciked %s\n", msg.Payload)
	}
}
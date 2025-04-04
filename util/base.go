package util

import "fmt"

func GetKey(uid string) string {
	return fmt.Sprintf("cart:{%s}", uid)
}
func GetKeySelect(uid string) string {
	return fmt.Sprintf("cart_select:{%s}", uid)
}
func GetKeyPrice(uid string) string {
	return fmt.Sprintf("cart_price:{%s}", uid)
}

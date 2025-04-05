namespace go eshop.cart

service cartService {
        // 添加商品
    BaseResponse addItem(1:AddItemRequest req),

        // 获取购物车列表
    PageResponse getList(1:PageRequest req),

        // 更新商品状态，返回总价
    UpdateResponse updateItem(1:UpdateRequest req),

        // 批量删除
    BaseResponse deleteItem(1:DeleteRequest req)
}
struct DeleteRequest {
    1: list<string> skus
    2: string uid
}
struct BaseResponse {
    2: i64 code
    3: optional string errStr
}
struct UpdateResponse {
    1: string price
    2: i64 code
    3: optional string errStr
}
struct PageRequest {
    1: i32 pageSize
    2: i32 pageNum
    3: string uid
}
struct AddItemRequest {
    1: required string skuId,
    2: required i32 quantity,
    3: string uid
}
struct UpdateRequest {
    1: optional i32 quantity,
    2: optional bool selected,
    3: required string skuId
    4: string uid
}
struct CartItem {
    1: string sku,
    2: i32 quantity,
    3: bool selected
}

struct PageResponse {
    1: i32 pageSize
    2: i32 pageNum
    3: bool isEnd
    4: list<CartItem> items
    5: optional string info
    6: string price
}

//kitex -module eshop_cart cart.thrift

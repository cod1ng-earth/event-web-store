module Message exposing (..)

import CatalogPage.Message
import CartPage.Message
import ProductDetailPage.Message


type Msg
    = CatalogPageMsg CatalogPage.Message.SubMsg
    | CartPageMsg CartPage.Message.SubMsg
    | ProductDetailPageMsg ProductDetailPage.Message.SubMsg
    | ShowCartPageMsg
    | ShowCatalogPage
    | ShowProductDetailPage String

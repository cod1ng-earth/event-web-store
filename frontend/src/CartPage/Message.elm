module CartPage.Message exposing (SubMsg(..))

import Http
import Checkout


type SubMsg
    = ChangeProductQuantity String Int
    | CartGotChanged (Result Http.Error Checkout.Cart)
    | OrderCart
    | LoadCart
    | CartGotOrdered (Result Http.Error Checkout.OrderCartResonse)

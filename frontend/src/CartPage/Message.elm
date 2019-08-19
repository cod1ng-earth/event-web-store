module CartPage.Message exposing (..)

import Http
import Checkout


type SubMsg
    = ChangeProductQuantity String Int
    | CartGotChanged (Result Http.Error Checkout.Cart)
    | OrderCart
    | CartGotOrdered (Result Http.Error Checkout.OrderCartResonse)

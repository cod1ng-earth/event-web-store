module CartPage.Model exposing (..)

import Checkout


type alias Model =
    { cart: MaybeOrderedCart
    , error : Maybe String
    }

type MaybeOrderedCart
    = Cart Checkout.Cart
    | OrderedCart

init : Model
init =
    { cart = Cart emptyCart
    , error = Nothing
    }

emptyCart = Checkout.Cart "" []

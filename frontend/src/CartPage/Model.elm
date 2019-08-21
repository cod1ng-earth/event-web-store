module CartPage.Model exposing (Model, MaybeOrderedCart(..), init)

import Task
import Time
import Checkout
import Message exposing (Msg)
import CartPage.Message


type alias Model =
    { cart : MaybeOrderedCart
    , error : Maybe String
    }


type MaybeOrderedCart
    = Cart Checkout.Cart
    | OrderedCart


init : ( Model, Cmd Msg )
init =
    let
        model =
            { cart = Cart emptyCart
            , error = Nothing
            }
    in
        ( model, Task.perform (\_ -> Message.CartPageMsg CartPage.Message.LoadCart) Time.here )


emptyCart : Checkout.Cart
emptyCart =
    Checkout.Cart "" []

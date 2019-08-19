module CartPage.Update exposing (..)

import Http
import Protobuf.Decode as Decode
import Protobuf.Encode as Encode
import Checkout
import CartPage.Model exposing (Model, MaybeOrderedCart(..), emptyCart)
import Message exposing (..)
import CartPage.Message exposing (..)


update : SubMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    case (msg, model.cart) of
        (ChangeProductQuantity uuid quantity, _) ->
            ( model, updateCart (Checkout.ChangeProductQuantity "" uuid quantity) )

        (OrderCart, _) ->
            ( model, orderCart )

        (CartGotChanged (Ok newCart), _) ->
            ( { model | cart = Cart newCart, error = Nothing }, Cmd.none )

        (CartGotChanged (Err e), _) ->
            ( { model | error = Just (toString e) }, Cmd.none )

        (CartGotOrdered (Ok orderState), _) ->
            if orderState.successful then                        
                ( { model | cart = OrderedCart, error = Nothing }, Cmd.none )
            else
                ( { model | error = Just "failed to order cart" }, Cmd.none ) 

        (CartGotOrdered (Err e), _) ->
            ( { model | error = Just (toString e) }, Cmd.none )               


toString : Http.Error -> String
toString error =
    case error of
        Http.BadUrl url ->
            url ++ " is bad"

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network Error"

        Http.BadStatus status ->
            String.fromInt status ++ " Status"

        Http.BadBody msg ->
            msg


updateCart : Checkout.ChangeProductQuantity -> Cmd Msg
updateCart cartChange =
    let
        e =
            Decode.expectBytes CartGotChanged Checkout.cartDecoder
    in
        Http.riskyRequest
            { method = "POST"
            , headers = []
            , url = "http://localhost:8080/cart"
            , body =
                Http.bytesBody "application/octet-stream" <|
                    Encode.encode (Checkout.toChangeProductQuantityEncoder cartChange)
            , expect = e
            , timeout = Nothing
            , tracker = Nothing
            }
            |> Cmd.map CartPageMsg


orderCart : Cmd Msg
orderCart =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = "http://localhost:8080/orderCart"
        , body = Http.emptyBody
        , expect = Decode.expectBytes CartGotOrdered Checkout.orderCartResonseDecoder
        , timeout = Nothing
        , tracker = Nothing
        }
        |> Cmd.map CartPageMsg


fetchCart : Cmd Msg
fetchCart =
    Http.riskyRequest
        { method = "GET"
        , headers = []
        , url = "http://localhost:8080/cart"
        , body = Http.emptyBody
        , expect = Decode.expectBytes CartGotChanged Checkout.cartDecoder
        , timeout = Nothing
        , tracker = Nothing
        }
        |> Cmd.map CartPageMsg

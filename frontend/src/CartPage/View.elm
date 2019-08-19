module CartPage.View exposing (..)

import Http
import Html exposing (Html, button, div, text, h1, h2, h3, ol, ul, li, a, span, header, footer, i, img, input)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (href, class, id, disabled, attribute, src, placeholder)
import Protobuf.Decode as Decode
import Protobuf.Encode as Encode
import List exposing (foldl, append)
import Checkout
import CatalogPage.View exposing (formatPrice, productImage)
import CartPage.Model exposing (..)
import CartPage.Update exposing (..)
import CartPage.Message exposing (..)
import Message exposing (..)


view : Model -> Html Msg
view model =
    case model.cart of
        Cart cart ->
            let
                error = case model.error of
                    Nothing ->
                        text ""
                    Just e ->
                        text e

                head =
                    List.map (\l -> renderCartItem l) cart.positions

                tail =
                    [ li [] [ button [ class "mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--accent", onClick (CartPageMsg OrderCart) ] [ text "Order Now" ] ] ]
            in
                div []
                [ error
                , ul [ class "product-list mdl-list" ] (append head tail)
                ]
        OrderedCart -> h2 [] [ text "Cart ordered sucessfully" ]

renderCartItem : Checkout.Position -> Html Msg
renderCartItem item =
    let
        quantity =
            item.quantity

        inStock =
            item.inStock

        outOfStock =
            not item.inStock

        moreInStock =
            item.moreInStock

        stockText =
            if not item.inStock then
                text "Out of Stock"
            else if not item.moreInStock then
                text "All stock in cart"
            else
                text ""
    in
        li [ class "mdl-list__item mdl-list__item--two-line" ]
            [ span [ class "mdl-list__item-primary-content" ]
                [ img [ class "custom-list-image", src (productImage item.productID 100 50) ] []

                --          , span [ onClick (LoadProduct item.productID) ] [ text item.title ]
                , span [ class "mdl-list__item-sub-title" ] [ text ("price: " ++ formatPrice item.price) ]
                ]
            , span []
                [ stockText
                , button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick (CartPageMsg (ChangeProductQuantity item.productID (quantity - 1))) ] [ i [ class "material-icons" ] [ text "remove" ] ]
                , span [] [ text (String.fromInt quantity) ]
                , button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick (CartPageMsg (ChangeProductQuantity item.productID (quantity + 1))), disabled (not moreInStock) ] [ i [ class "material-icons" ] [ text "add" ] ]
                ]
            ]


itemsInCart : MaybeOrderedCart -> String
itemsInCart maybeOrderedCart =
    case maybeOrderedCart of
        Cart cart ->
            let
                inStock =
                    itemsInCartInStock cart

                outOfStock =
                    itemsInCartOutOfStock cart
            in
                if outOfStock == 0 then
                    String.fromInt inStock
                else
                    String.fromInt (inStock + outOfStock) ++ "-" ++ String.fromInt outOfStock
        OrderedCart ->
            "0"


itemsInCartInStock : Checkout.Cart -> Int
itemsInCartInStock cart =
    let
        ls =
            List.map
                (\e ->
                    if e.inStock then
                        1
                    else
                        0
                )
                cart.positions
    in
        foldl (+) 0 ls


itemsInCartOutOfStock : Checkout.Cart -> Int
itemsInCartOutOfStock cart =
    let
        ls =
            List.map
                (\e ->
                    if not e.inStock then
                        1
                    else
                        0
                )
                cart.positions
    in
        foldl (+) 0 ls

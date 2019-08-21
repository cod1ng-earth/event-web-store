module View exposing (..)

import Html exposing (button, div, h1, header, Html, span, text)
import Html.Attributes exposing (attribute, class, disabled, id)
import Html.Events exposing (onClick)

import Model exposing (..)
import Message exposing (..)

import CatalogPage.View
import CartPage.View exposing (itemsInCart)
import CartPage.Message exposing (..)
import ProductDetailPage.View


view : Model -> Html Msg
view model =
    let
        ( showingCart, showingCatalog ) =
            case model.content of
                CatalogPage _ ->
                    ( False, True )

                ProductDetailPage _ ->
                    ( False, False )

                ShowCartPage ->
                    ( True, False )

                OrderSuccessfulPage ->
                    ( False, False )

                OrderFailedPage ->
                    ( False, False )
    in
        div [ class "mdl-layout mdl-layout--fixed-header" ]
            [ header [ class "mdl-layout__header mdl-layout__header--waterfall custom-header" ]
                [ div [ class "mdl-layout__header-row custom-header-row" ]
                    [ span [ class "mdl-layout__title" ] [ text "Event Thingy Store" ]
                    , div [ class "mdl-layout-spacer" ] []
                    , div [] [ span [ class "mdl-badge custom-header-cart", attribute "data-badge" (itemsInCart model.cart.cart), onClick ShowCartPageMsg ] [ text "Cart" ] ]
                    , button [ class "mdl-button mdl-button--raised mdl-button--accent", onClick ShowCartPageMsg, disabled showingCart ] [ text "show Cart" ]
                    , button [ class "mdl-button mdl-button--raised mdl-button--accent", onClick ShowCatalogPage, disabled showingCatalog ] [ text "show products" ]
                    ]
                ]
            , div [ class "custom-header-error" ] [ text model.error ]
            , div [ class "mdl-layout__content", id "main" ] [ renderContent model ]
            ]


renderContent : Model -> Html Msg
renderContent model =
    case model.content of
        CatalogPage mdl ->
            CatalogPage.View.view mdl

        ProductDetailPage mdl ->
            ProductDetailPage.View.view mdl

        ShowCartPage ->
            CartPage.View.view model.cart

        OrderSuccessfulPage ->
            h1 [] [ text "Order was successful" ]

        OrderFailedPage ->
            h1 [] [ text "Ordering failed" ]

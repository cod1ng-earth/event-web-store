module CatalogPage.View exposing (..)

import Html exposing (a, button, div, Html, i, img, input, li, span, text, ul)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (class, disabled, placeholder, src)
import Char exposing (toCode)
import String exposing (toList)
import Catalog
import CatalogPage.Model exposing (Model)
import CatalogPage.Update exposing (..)
import CatalogPage.Message exposing (..)
import ProductDetailPage.Message exposing (..)
import CartPage.Message exposing (..)
import Message exposing (..)


view : Model -> Html Msg
view cp =
    let
        sorting =
            cp.sorting

        filtering =
            cp.filtering

        prevEnabled =
            cp.currentPage > 0

        nextEnabled =
            cp.currentPage < (cp.totalPages - 1)

        uuidDisabled =
            sorting == "uuid"

        priceDisabled =
            sorting == "price"

        nameDisabled =
            sorting == "name"

        prefixDisabled =
            filtering == "prefix"

        pagesText =
            if cp.totalPages == 0 then
                text " No luck! "
            else
                text (" Page " ++ String.fromInt (cp.currentPage + 1) ++ " from " ++ String.fromInt (cp.totalPages))

        pp =
            case cp.products of
                Nothing ->
                    []

                Just products ->
                    List.map (\l -> renderProduct l) products
    in
        div [ class "mdl-grid" ]
            [ div [ class "mdl-cell mdl-cell--12-col" ]
                [ button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick (CatalogPageMsg PreviousPage), disabled (not prevEnabled) ] [ i [ class "material-icons" ] [ text "remove" ] ]
                , span [ class "mdl-title custom-page-number" ]
                    [ pagesText
                    , input [ class "custom-go-to-page", placeholder "Go to page", onInput pageLoader ] []
                    ]
                , button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick (CatalogPageMsg NextPage), disabled (not nextEnabled) ] [ i [ class "material-icons" ] [ text "add" ] ]
                , span [ class "custom-sorting" ]
                    [ text "Sort by "
                    , sortProductsButton SortByName nameDisabled "Name"
                    , sortProductsButton SortByUuid uuidDisabled "Uuid"
                    , sortProductsButton SortByPrice priceDisabled "Price"
                    ]
                , span [ class "custom-sorting" ]
                    [ text "Filter by "
                    , filterProductsButton prefixDisabled "Prefix"
                    , input [ onInput prefixFilter, disabled (not prefixDisabled) ] []
                    ]
                ]
            , ul [ class "product-list mdl-list" ] pp
            ]


pageLoader : String -> Message.Msg
pageLoader str =
    CatalogPageMsg (GoToPage str)


prefixFilter : String -> Message.Msg
prefixFilter str =
    CatalogPageMsg (SetFilterPrefix str)


sortProductsButton : CatalogPage.Message.SubMsg -> Bool -> String -> Html Msg
sortProductsButton click inactive label =
    button
        [ class "mdl-button mdl-js-button mdl-button--raised mdl-button--colored", onClick (CatalogPageMsg click), disabled inactive ]
        [ text label ]


filterProductsButton : Bool -> String -> Html Msg
filterProductsButton active label =
    if active then
        button
            [ class "mdl-button mdl-button--raised mdl-button--accent", onClick (CatalogPageMsg DisableFilterByPrefix) ]
            [ text label ]
    else
        button
            [ class "mdl-button mdl-js-button mdl-button--raised mdl-button--colored", onClick (CatalogPageMsg EnableFilterByPrefix) ]
            [ text label ]


renderProduct : Catalog.Product -> Html Msg
renderProduct product =
    li [ class "mdl-list__item mdl-list__item--two-line" ]
        [ span [ class "mdl-list__item-primary-content" ]
            [ img [ class "custom-list-image", src (productImage product.id 100 50) ] []
            , span [] [ text product.name ]
            , span [ onClick (ProductDetailPageMsg (LoadProduct product.id)) ] [ text product.name ]
            , span [ class "mdl-list__item-sub-title" ] [ text ("price: " ++ formatPrice product.price) ]
            ]
        , span [ class "mdl-list__item-secondary-content" ] [ showProductButton product.id ]
        , span [ class "mdl-list__item-secondary-content" ] [ addToCartButton product.id ]
        ]


productImage : String -> Int -> Int -> String
productImage uuid width height =
    "https://picsum.photos/id/" ++ String.fromInt (modBy 50 (reduceUuid uuid)) ++ "/" ++ String.fromInt width ++ "/" ++ String.fromInt height


reduceUuid : String -> Int
reduceUuid uuid =
    List.foldl (\x a -> x + a) 0 (List.map toCode (toList uuid))


formatPrice : Int -> String
formatPrice price =
    String.fromInt (price // 100) ++ "€"


showProductButton : String -> Html Msg
showProductButton productID =
    button
        [ class "mdl-button mdl-js-button mdl-button--raised mdl-button--colored", onClick (ShowProductDetailPage productID) ]
        [ text "show Details" ]


addToCartButton : String -> Html Msg
addToCartButton uuid =
    button
        [ class "mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--accent", onClick (CartPageMsg (ChangeProductQuantity uuid 1)) ]
        [ text "add to cart" ]

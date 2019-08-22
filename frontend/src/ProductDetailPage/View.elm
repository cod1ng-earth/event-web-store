module ProductDetailPage.View exposing (..)

import Html exposing (div, h1, Html, img, span, text)
import Html.Attributes exposing (class, src)
import CatalogPage.View exposing (addToCartButton, formatPrice, productImage)
import ProductDetailPage.Model exposing (..)
import ProductDetailPage.Update exposing (..)
import ProductDetailPage.Message exposing (..)
import Message exposing (..)


view : Model -> Html Msg
view model =
    case model of
        Loading ->
            text ""

        LoadingSlowly ->
            text "loading..."

        Loaded product ->
            div [ class "mdl-grid" ]
                [ div [ class "mdl-cell mdl-cell--12-col" ] [ h1 [] [ text product.name ] ]
                , div [ class "mdl-cell mdl-cell--8-col" ] [ img [ class "custom-detail-image", src (productImage product.id 400 200) ] [] ]
                , div [ class "mdl-cell mdl-cell--4-col" ]
                    [ span [ class "custom-detail-block" ] [ text ("id: " ++ product.id) ]
                    , span [ class "mdl-typography--headline custom-detail-block" ] [ text product.description ]
                    , span [ class "custom-detail-block" ] [ text product.longtext ]
                    , span [ class "mdl-typography--display-1 custom-detail-block" ] [ text (formatPrice product.price) ]
                    , span [] [ addToCartButton product.id ]
                    ]
                ]

        Failed err ->
            text err

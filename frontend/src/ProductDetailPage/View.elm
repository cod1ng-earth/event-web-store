module ProductDetailPage.View exposing (..)

import Html exposing (Html, button, div, text, h1, h2, h3, ol, ul, li, a, span, header, footer, i, img, input)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (href, class, id, disabled, attribute, src, placeholder)
import Http
import Protobuf.Decode as Decode
import Catalog
import Checkout
import CatalogPage.View exposing (productImage, formatPrice, addToCartButton)
import ProductDetailPage.Model exposing (Model)
import ProductDetailPage.Update exposing (..)
import ProductDetailPage.Message exposing (..)
import Message exposing (..)


view : Model -> Html Msg
view model =
    let
        maybeProduct =
            model.product
    in
        case maybeProduct of
            Nothing ->
                text "sorry no product"

            Just product ->
                div [ class "mdl-grid" ]
                    [ div [ class "mdl-cell mdl-cell--12-col" ] [ h1 [] [ text product.title ] ]
                    , div [ class "mdl-cell mdl-cell--8-col" ] [ img [ class "custom-detail-image", src (productImage product.id 400 200) ] [] ]
                    , div [ class "mdl-cell mdl-cell--4-col" ]
                        [ span [ class "custom-detail-block" ] [ text ("id: " ++ product.id) ]
                        , span [ class "mdl-typography--headline custom-detail-block" ] [ text product.description ]
                        , span [ class "custom-detail-block" ] [ text product.longtext ]
                        , span [ class "mdl-typography--display-1 custom-detail-block" ] [ text (formatPrice product.price) ]
                        , span [] [ addToCartButton product.id ]
                        ]
                    ]

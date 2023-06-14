// Code generated by go-localize; DO NOT EDIT.
// This file was generated by robots at
// 2023-06-13 18:01:15.435034 +0100 IST m=+0.004341585

package localizations

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var localizations = map[string]string{
	"en.messages.added-required-request-body":                              "added required request body",
	"en.messages.api-deprecated-sunset-parse":                              "api sunset date '%s' can't be parsed for deprecated API: %v",
	"en.messages.api-operation-id-added":                                   "api operation id %s was added",
	"en.messages.api-operation-id-removed":                                 "api operation id %s removed and replaced with %s",
	"en.messages.api-path-removed-before-sunset":                           "api path removed before the sunset date %s",
	"en.messages.api-path-removed-without-deprecation":                     "api path removed without deprecation",
	"en.messages.api-removed-before-sunset":                                "api removed before the sunset date %s",
	"en.messages.api-removed-without-deprecation":                          "api removed without deprecation",
	"en.messages.api-schema-removed":                                       "removed the schema %s from openapi components",
	"en.messages.api-sunset-date-changed-too-small":                        "api sunset date changed to earlier date from %s to %s, new sunset date must be not earlier than %s at least %d days from now",
	"en.messages.api-sunset-date-too-small":                                "api sunset date '%s' is too small, must be at least %d days from now",
	"en.messages.api-tag-added":                                            "api tag %s added",
	"en.messages.api-tag-removed":                                          "api tag %s removed",
	"en.messages.at":                                                       "at",
	"en.messages.endpoint-added":                                           "endpoint added",
	"en.messages.endpoint-deprecated":                                      "endpoint deprecated",
	"en.messages.endpoint-reactivated":                                     "endpoint reactivated",
	"en.messages.in":                                                       "in",
	"en.messages.new-optional-request-parameter":                           "added the new optional %s request parameter %s",
	"en.messages.new-request-path-parameter":                               "added the new path request parameter %s",
	"en.messages.new-required-request-header-property":                     "added the new required %s request header's property %s",
	"en.messages.new-required-request-parameter":                           "added the new required %s request parameter %s",
	"en.messages.new-required-request-property":                            "added the new required request property %s",
	"en.messages.optional-response-header-removed":                         "the optional response header %s removed for the status %s",
	"en.messages.pattern-changed-warn-comment":                             "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern(e.g. changed from '[0-9]+' to '[0-9]*')",
	"en.messages.request-allOf-modified":                                   "modified allOf for the request property %s",
	"en.messages.request-allOf-modified-comment":                           "It is a warning because it is very difficult to check that allOf changed correctly without breaking changes",
	"en.messages.request-body-became-enum":                                 "request body was restricted to a list of enum values",
	"en.messages.request-body-became-not-nullable":                         "the request's body became not nullable",
	"en.messages.request-body-became-optional":                             "request body became optional",
	"en.messages.request-body-became-required":                             "request body became required",
	"en.messages.request-body-enum-value-removed":                          "request body enum value removed %s",
	"en.messages.request-body-max-decreased":                               "the request's body max was decreased to %s",
	"en.messages.request-body-max-length-decreased":                        "the request's body maxLength was decreased to %s",
	"en.messages.request-body-max-length-set":                              "the request's body maxLength was set to %s",
	"en.messages.request-body-max-length-set-comment":                      "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-body-max-set":                                     "the request's body max was set to %s",
	"en.messages.request-body-max-set-comment":                             "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-body-min-increased":                               "the request's body min was increased to %s",
	"en.messages.request-body-min-items-increased":                         "the request's body minItems was increased to %s",
	"en.messages.request-body-min-items-set":                               "the request's body minItems was set to %s",
	"en.messages.request-body-min-items-set-comment":                       "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-body-min-set":                                     "the request's body min was set to %s",
	"en.messages.request-body-min-set-comment":                             "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-body-type-changed":                                "the request's body type/format changed from %s/%s to %s/%s",
	"en.messages.request-header-property-became-enum":                      "the %s request header's property %s was restricted to a list of enum values",
	"en.messages.request-header-property-became-required":                  "the %s request header's property %s became required",
	"en.messages.request-parameter-became-enum":                            "the %s request parameter %s was restricted to a list of enum values",
	"en.messages.request-parameter-became-optional":                        "the %s request parameter %s became optional",
	"en.messages.request-parameter-became-required":                        "the %s request parameter %s became required",
	"en.messages.request-parameter-default-value-changed":                  "for the %s request parameter %s, default value was changed from %s to %s",
	"en.messages.request-parameter-enum-value-removed":                     "removed the enum value %s for the %s request parameter %s",
	"en.messages.request-parameter-max-decreased":                          "for the %s request parameter %s, the max was decreased from %s to %s",
	"en.messages.request-parameter-max-length-decreased":                   "for the %s request parameter %s, the maxLength was decreased from %s to %s",
	"en.messages.request-parameter-max-length-set":                         "for the %s request parameter %s, the maxLength was set to %s",
	"en.messages.request-parameter-max-length-set-comment":                 "This is a warning because sometimes it is required to be set because of security reasons or current error in specification. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-parameter-max-set":                                "for the %s request parameter %s, the max was set to %s",
	"en.messages.request-parameter-max-set-comment":                        "This is a warning because sometimes it is required to be set because of security reasons or current error in specification. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-parameter-min-increased":                          "for the %s request parameter %s, the min was increased from %s to %s",
	"en.messages.request-parameter-min-items-increased":                    "for the %s request parameter %s, the minItems was increased from %s to %s",
	"en.messages.request-parameter-min-items-set":                          "for the %s request parameter %s, the minItems was set to %s",
	"en.messages.request-parameter-min-items-set-comment":                  "This is a warning because sometimes it is required to be set because of security reasons or current error in specification. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-parameter-min-set":                                "for the %s request parameter %s, the min was set to %s",
	"en.messages.request-parameter-min-set-comment":                        "This is a warning because sometimes it is required to be set because of security reasons or current error in specification. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-parameter-pattern-added":                          "added the pattern '%s' for the %s request parameter %s",
	"en.messages.request-parameter-pattern-changed":                        "changed the pattern for the %s request parameter %s from '%s' to '%s'",
	"en.messages.request-parameter-removed":                                "deleted the %s request parameter %s",
	"en.messages.request-parameter-type-changed":                           "for the %s request parameter %s, the type/format was changed from %s/%s to %s/%s",
	"en.messages.request-parameter-x-extensible-enum-value-removed":        "removed the x-extensible-enum value %s for the %s request parameter %s",
	"en.messages.request-property-became-enum":                             "request property %s was restricted to a list of enum values",
	"en.messages.request-property-became-not-nullable":                     "the request property %s became not nullable",
	"en.messages.request-property-became-required":                         "the request property %s became required",
	"en.messages.request-property-enum-value-removed":                      "removed the enum value %s of the request property %s",
	"en.messages.request-property-max-decreased":                           "the %s request property's max was decreased to %s",
	"en.messages.request-property-max-length-decreased":                    "the %s request property's maxLength was decreased to %s",
	"en.messages.request-property-max-length-set":                          "the %s request property's maxLength was set to %s",
	"en.messages.request-property-max-length-set-comment":                  "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-property-max-set":                                 "the %s request property's max was set to %s",
	"en.messages.request-property-max-set-comment":                         "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-property-min-increased":                           "the %s request property's min was increased to %s",
	"en.messages.request-property-min-items-increased":                     "the %s request property's minItems was increased to %s",
	"en.messages.request-property-min-items-set":                           "the %s request property's minItems was set to %s",
	"en.messages.request-property-min-items-set-comment":                   "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-property-min-set":                                 "the %s request property's min was set to %s",
	"en.messages.request-property-min-set-comment":                         "This is a warning because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
	"en.messages.request-property-pattern-added":                           "added the pattern '%s' for the request property %s",
	"en.messages.request-property-pattern-changed":                         "changed the pattern for the request property %s from '%s' to '%s'",
	"en.messages.request-property-removed":                                 "removed the request property %s",
	"en.messages.request-property-type-changed":                            "the %s request property type/format changed from %s/%s to %s/%s",
	"en.messages.request-property-x-extensible-enum-value-removed":         "removed the x-extensible-enum value '%s' of the request property %s",
	"en.messages.required-response-header-removed":                         "the mandatory response header %s removed for the status %s",
	"en.messages.response-body-became-nullable":                            "the response's body became nullable",
	"en.messages.response-body-max-increased":                              "the response's body max was increased from %s to %s",
	"en.messages.response-body-max-length-increased":                       "the response's body maxLength was increased from %s to %s",
	"en.messages.response-body-max-length-unset":                           "the response's body maxLength was unset from %s",
	"en.messages.response-body-min-decreased":                              "the response's body min was decreased from %s to %s",
	"en.messages.response-body-min-items-decreased":                        "the response's body minItems was decreased from %s to %s",
	"en.messages.response-body-min-items-unset":                            "the response's body minItems was unset from %s",
	"en.messages.response-body-min-length-decreased":                       "the response's body minLength was decreased from %s to %s",
	"en.messages.response-body-type-changed":                               "the response's body type/format changed from %s/%s to %s/%s for status %s",
	"en.messages.response-header-became-optional":                          "the response header %s became optional for the status %s",
	"en.messages.response-media-type-removed":                              "removed the media type %s for the response with the status %s",
	"en.messages.response-mediatype-enum-value-removed":                    "response schema %s enum value removed %s",
	"en.messages.response-non-success-status-added":                        "added the non-success response with the status %s",
	"en.messages.response-non-success-status-removed":                      "removed the non-success response with the status %s",
	"en.messages.response-optional-property-removed":                       "removed the optional property %s from the response with the %s status",
	"en.messages.response-property-became-nullable":                        "the response property %s became nullable for the status %s",
	"en.messages.response-property-became-optional":                        "the response property %s became optional for the status %s",
	"en.messages.response-property-enum-value-added":                       "added the new '%s' enum value the %s response property for the response status %s",
	"en.messages.response-property-enum-value-added-comment":               "Adding new enum values to response could be unexpected for clients, use x-extensible-enum instead.",
	"en.messages.response-property-enum-value-removed":                     "removed the '%s' enum value from the %s response property for the response status %s",
	"en.messages.response-property-max-increased":                          "the %s response property's max was increased from %s to %s for the response status %s",
	"en.messages.response-property-max-length-increased":                   "the %s response property's maxLength was increased from %s to %s for the response status %s",
	"en.messages.response-property-max-length-unset":                       "the %s response property's maxLength was unset from %s for the response status %s",
	"en.messages.response-property-min-decreased":                          "the %s response property's min was decreased from %s to %s for the response status %s",
	"en.messages.response-property-min-items-decreased":                    "the %s response property's minItems was decreased from %s to %s for the response status %s",
	"en.messages.response-property-min-items-unset":                        "the %s response property's minItems was unset from %s for the response status %s",
	"en.messages.response-property-min-length-decreased":                   "the %s response property's minLength was decreased from %s to %s for the response status %s",
	"en.messages.response-property-type-changed":                           "the response's property type/format changed from %s/%s to %s/%s for status %s",
	"en.messages.response-required-property-became-not-write-only":         "the response required property %s became not write-only for the status %s",
	"en.messages.response-required-property-became-not-write-only-comment": "It is valid only if the property was always returned before the specification has been changed",
	"en.messages.response-required-property-removed":                       "removed the required property %s from the response with the %s status",
	"en.messages.response-success-status-added":                            "added the success response with the status %s",
	"en.messages.response-success-status-removed":                          "removed the success response with the status %s",
	"en.messages.sunset-deleted":                                           "api sunset date deleted, but deprecated=true kept",
	"en.messages.total-errors":                                             "Backward compatibility errors (%d):\n",
	"ru.messages.added-required-request-body":                              "добавлено обязательное тело запроса",
	"ru.messages.api-deprecated-sunset-parse":                              "API deprecated без валидно парсящейся '%s' даты sunset: %v",
	"ru.messages.api-operation-id-added":                                   "добавлен идентификатор операции API %s",
	"ru.messages.api-operation-id-removed":                                 "Идентификатор операции API %s удален и заменен на %s",
	"ru.messages.api-path-added":                                           "API path добавлено",
	"ru.messages.api-path-deprecated":                                      "API path deprecated",
	"ru.messages.api-path-reactivated":                                     "API path реактивирован",
	"ru.messages.api-path-removed-before-sunset":                           "API path удалён до даты sunset %s",
	"ru.messages.api-path-removed-without-deprecation":                     "api path удалён без процедуры deprecation",
	"ru.messages.api-removed-before-sunset":                                "API удалёг до даты sunset %s",
	"ru.messages.api-removed-without-deprecation":                          "API удалён без deprecation",
	"ru.messages.api-schema-removed":                                       "удалена схема %s из компонентов openapi",
	"ru.messages.api-sunset-date-changed-too-small":                        "дата sunset у API изменена на более раннюю с %s на %s, новая дата sunset должна быть либо не раньше %s, либо, как минимум, %d дней от текущего дня",
	"ru.messages.api-sunset-date-too-small":                                "дата API sunset date '%s' слишком ранняя, должно быть как минимум %d дней от текущего дня",
	"ru.messages.api-tag-added":                                            "тег API %s добавлен",
	"ru.messages.api-tag-removed":                                          "Тег API %s удален",
	"ru.messages.at":                                                       "в",
	"ru.messages.in":                                                       "в",
	"ru.messages.new-optional-request-parameter":                           "добавлен новый необязательный %s параметр зароса %s",
	"ru.messages.new-request-path-parameter":                               "добален новый path параметр запроса %s",
	"ru.messages.new-required-request-header-property":                     "в заголовке запроса %s добавлено новое обязательное поле %s",
	"ru.messages.new-required-request-parameter":                           "добавлен новый обязательный %s параметр зароса %s",
	"ru.messages.new-required-request-property":                            "добавлено новле обязательное поле запроса %s",
	"ru.messages.optional-response-header-removed":                         "удалён ранее необязательный заголовок ответа %s для ответа со статусом %s",
	"ru.messages.pattern-changed-warn-comment":                             "Это предупреждение, потому что сложно автоматически проанализировать, является ли новый шаблон надмножеством предыдущего шаблона (например, изменен с '[0-9]+' на '[0-9]*').",
	"ru.messages.request-allOf-modified":                                   "изменено allOf для поля запроса %s",
	"ru.messages.request-allOf-modified-comment":                           "Это предупреждение, потому что очень сложно алгоритмически автоматизированно проверить правильность изменения allOf на обратную совместимость.",
	"ru.messages.request-body-became-enum":                                 "тело запроса было ограничено списком значений перечисления",
	"ru.messages.request-body-became-not-nullable":                         "тело запроса стало недействительным",
	"ru.messages.request-body-became-optional":                             "тело запроса стало необязательным",
	"ru.messages.request-body-became-required":                             "тело запроса стало обязательным",
	"ru.messages.request-body-enum-value-removed":                          "значение перечисления тела запроса удалено %s",
	"ru.messages.request-body-max-decreased":                               "значение max у тела запроса уменьшено до %s",
	"ru.messages.request-body-max-length-decreased":                        "значение maxLength у тела запроса уменьшено до %s",
	"ru.messages.request-body-max-length-set":                              "у тела запроса задано значение maxLength в %s",
	"ru.messages.request-body-max-length-set-comment":                      "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-body-max-set":                                     "у тела запроса задано значение max в %s",
	"ru.messages.request-body-max-set-comment":                             "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-body-min-increased":                               "значение min у тела запроса увеличено до %s",
	"ru.messages.request-body-min-items-increased":                         "значение minItems у тела запроса увеличено до %s",
	"ru.messages.request-body-min-items-set":                               "задано значение minItems у тела запроса в %s",
	"ru.messages.request-body-min-items-set-comment":                       "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-body-min-set":                                     "задано значение min у тела запроса в %s",
	"ru.messages.request-body-min-set-comment":                             "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-body-type-changed":                                "изменился type/format тела запроса с %s/%s на %s/%s",
	"ru.messages.request-header-property-became-enum":                      "свойство %s заголовка запроса %s было ограничено списком значений перечисления",
	"ru.messages.request-header-property-became-required":                  "в заголовке запроса %s поле %s стало обязательным",
	"ru.messages.request-parameter-became-enum":                            "заголовок запроса %s поле %s было ограничено списком значений перечисления",
	"ru.messages.request-parameter-became-optional":                        "ранее необязательный параметр запроса %s %s теперь является необязательным",
	"ru.messages.request-parameter-became-required":                        "ранее необязательный %s параметр запроса %s стал обязательным",
	"ru.messages.request-parameter-default-value-changed":                  "в %s параметре запроса %s, значение по умолчанию изменено с %s на %s",
	"ru.messages.request-parameter-enum-value-removed":                     "удалено значение enum %s у %s параметра запроса %s",
	"ru.messages.request-parameter-max-decreased":                          "в %s параметре запроса %s, max уменьшен с %s до %s",
	"ru.messages.request-parameter-max-length-decreased":                   "в %s параметре запроса %s, maxLength уменьшен с %s до %s",
	"ru.messages.request-parameter-max-length-set":                         "в %s параметре запроса %s, maxLength установлен в %s",
	"ru.messages.request-parameter-max-length-set-comment":                 "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-parameter-max-set":                                "в %s параметре запроса %s, max установлен в %s",
	"ru.messages.request-parameter-max-set-comment":                        "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-parameter-min-increased":                          "в %s параметре запроса %s, min увеличен с %s до %s",
	"ru.messages.request-parameter-min-items-increased":                    "в %s параметре запроса %s, minItems увеличен с %s до %s",
	"ru.messages.request-parameter-min-items-set":                          "в %s параметре запроса %s, minItems установлен в %s",
	"ru.messages.request-parameter-min-items-set-comment":                  "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-parameter-min-set":                                "в %s параметре запроса %s, min установлен в %s",
	"ru.messages.request-parameter-min-set-comment":                        "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-parameter-pattern-added":                          "добавлен pattern '%s' у %s параметра запроса %s",
	"ru.messages.request-parameter-pattern-changed":                        "изменён pattern у %s параметра запроса %s со значения '%s' на значение '%s'",
	"ru.messages.request-parameter-removed":                                "удалён %s параметр запроса %s",
	"ru.messages.request-parameter-type-changed":                           "в %s параметре запроса %s, type/format изменился с %s/%s на %s/%s",
	"ru.messages.request-parameter-x-extensible-enum-value-removed":        "удалено из x-extensible-enum значение %s у %s параметра запроса %s",
	"ru.messages.request-property-became-enum":                             "свойство запроса %s было ограничено списком значений перечисления",
	"ru.messages.request-property-became-not-nullable":                     "свойство запроса %s стало недействительным",
	"ru.messages.request-property-became-required":                         "поле запроса %s стало обязательным",
	"ru.messages.request-property-enum-value-removed":                      "удалено enum значение %s у поля запроса %s",
	"ru.messages.request-property-max-decreased":                           "значение max у поля запроса %s уменьшено до %s",
	"ru.messages.request-property-max-length-decreased":                    "значение maxLength у поля запроса %s уменьшено до %s",
	"ru.messages.request-property-max-length-set":                          "у поля запроса %s задано значение maxLength в %s",
	"ru.messages.request-property-max-length-set-comment":                  "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-property-max-set":                                 "у поля запроса %s задано значение max в %s",
	"ru.messages.request-property-max-set-comment":                         "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-property-min-increased":                           "у поля запроса %s, увеличено значение min до %s",
	"ru.messages.request-property-min-items-increased":                     "значение minItems у поля запроса %s увеличено до %s",
	"ru.messages.request-property-min-items-set":                           "у поля запроса %s задано значение minItems в %s",
	"ru.messages.request-property-min-items-set-comment":                   "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-property-min-set":                                 "у поля запроса %s задано значение min в %s",
	"ru.messages.request-property-min-set-comment":                         "Это предупреждение, потому что иногда его требуется установить из соображений безопасности или из-за текущей ошибки в спецификации. Но хорошие клиенты должны быть проверены на поддержку этого ограничения перед внесением таких изменений в спецификацию.",
	"ru.messages.request-property-pattern-added":                           "добавлен pattern '%s' у поля запроса %s",
	"ru.messages.request-property-pattern-changed":                         "изменён pattern у поля запроса %s со значения '%s' на значение '%s'",
	"ru.messages.request-property-removed":                                 "удалено поле запроса %s",
	"ru.messages.request-property-type-changed":                            "у поля запроса %s изменился type/format с %s/%s на %s/%s",
	"ru.messages.request-property-x-extensible-enum-value-removed":         "удалено значение x-extensible-enum '%s' в поле запроса %s",
	"ru.messages.required-response-header-removed":                         "удалён ранее обязательный заголовок ответа %s для ответа со статусом %s",
	"ru.messages.response-body-became-nullable":                            "у тела ответа стало обнуляемым",
	"ru.messages.response-body-max-increased":                              "у тела ответа max увеличен с %s до %s",
	"ru.messages.response-body-max-length-increased":                       "у тела ответа maxLength увеличен с %s до %s",
	"ru.messages.response-body-max-length-unset":                           "у тела ответа maxLength был удалён, предыдущее значение - %s",
	"ru.messages.response-body-min-decreased":                              "у тела ответа min уменьшено с %s до %s",
	"ru.messages.response-body-min-items-decreased":                        "у тела ответа minItems уменьшено с %s до %s",
	"ru.messages.response-body-min-items-unset":                            "удалено значение minItems для тела ответа, предыдущее значение - %s",
	"ru.messages.response-body-min-length-decreased":                       "значение minLength для тела ответа уменьшено с %s до %s",
	"ru.messages.response-body-type-changed":                               "у тела ответа type/format изменился с %s/%s на %s/%s для ответа со статусом %s",
	"ru.messages.response-header-became-optional":                          "заголовок ответа %s стал необязательным для ответа со статусом %s",
	"ru.messages.response-media-type-removed":                              "удалён media type %s для ответа со статусом %s",
	"ru.messages.response-mediatype-enum-value-removed":                    "значение перечисления схемы ответа %s удалено %s",
	"ru.messages.response-non-success-status-added":                        "добавлен ответ об отсутствии успеха со статусом %s",
	"ru.messages.response-non-success-status-removed":                      "удален неуспешный (не 2xx) статус ответа %s",
	"ru.messages.response-optional-property-removed":                       "удалено необязательное поле %s из ответа со статусом %s",
	"ru.messages.response-property-became-nullable":                        "поле ответа %s стало обнуляемым для ответа со статусом %s",
	"ru.messages.response-property-became-optional":                        "поле ответа %s стало необязательным для ответа со статусом %s",
	"ru.messages.response-property-enum-value-added":                       "добавлено новое enum значение %s в поле ответа %s для ответа со статусом %s",
	"ru.messages.response-property-enum-value-added-comment":               "Добавление новых значений перечисления в ответ может быть неожиданным для клиентов, вместо этого используйте x-extensible-enum.",
	"ru.messages.response-property-enum-value-removed":                     "удалено значение перечисления '%s' из свойства ответа %s для статуса ответа %s.",
	"ru.messages.response-property-max-increased":                          "у поля ответа %s max увеличен с %s до %s для ответа со статусом %s",
	"ru.messages.response-property-max-length-increased":                   "у поля ответа %s maxLength увеличен с %s до %s для ответа со статусом %s",
	"ru.messages.response-property-max-length-unset":                       "у поля ответа %s maxLength был удалён, предыдущее значение - %s, для ответа со статусом %s",
	"ru.messages.response-property-min-decreased":                          "для поля ответа %s min уменьшен с %s до %s для ответа со статусом %s",
	"ru.messages.response-property-min-items-decreased":                    "у поля ответа %s уменьшено minItems с %s до %s для ответа со статусом %s",
	"ru.messages.response-property-min-items-unset":                        "у поля ответа %s удалено значение minItems, предыдущее значение - %s, для ответа со статусом %s",
	"ru.messages.response-property-min-length-decreased":                   "для поля ответа %s minLength уменьшен с %s до %s для ответа со статусом %s",
	"ru.messages.response-property-type-changed":                           "у поля type/format изменился с %s/%s на %s/%s для ответа со статусом %s",
	"ru.messages.response-required-property-became-not-write-only":         "обязательное поле ответа %s перестало быть write-only для ответа со статусом %s",
	"ru.messages.response-required-property-became-not-write-only-comment": "Изменение допустимо только в том случае, если свойство ВСЕГДА возвращалось ДО изменения спецификации.",
	"ru.messages.response-required-property-removed":                       "удалено обязательное поле ответа %s из ответа со статусом %s",
	"ru.messages.response-success-status-added":                            "добавлен ответ об успехе со статусом %s",
	"ru.messages.response-success-status-removed":                          "удален успешный (2xx) статус ответа %s",
	"ru.messages.sunset-deleted":                                           "удалена дата sunset date у API, но сохранён deprecated=true",
	"ru.messages.total-errors":                                             "Ошибки обратной совместимости (всего: %d):\n",
}

type Replacements map[string]interface{}

type Localizer struct {
	Locale         string
	FallbackLocale string
	Localizations  map[string]string
}

func New(locale string, fallbackLocale string) *Localizer {
	t := &Localizer{Locale: locale, FallbackLocale: fallbackLocale}
	t.Localizations = localizations
	return t
}

func (t Localizer) SetLocales(locale, fallback string) Localizer {
	t.Locale = locale
	t.FallbackLocale = fallback
	return t
}

func (t Localizer) SetLocale(locale string) Localizer {
	t.Locale = locale
	return t
}

func (t Localizer) SetFallbackLocale(fallback string) Localizer {
	t.FallbackLocale = fallback
	return t
}

func (t Localizer) GetWithLocale(locale, key string, replacements ...*Replacements) string {
	str, ok := t.Localizations[t.getLocalizationKey(locale, key)]
	if !ok {
		str, ok = t.Localizations[t.getLocalizationKey(t.FallbackLocale, key)]
		if !ok {
			return key
		}
	}

	// If the str doesn't have any substitutions, no need to
	// template.Execute.
	if strings.Index(str, "}}") == -1 {
		return str
	}

	return t.replace(str, replacements...)
}

func (t Localizer) Get(key string, replacements ...*Replacements) string {
	str := t.GetWithLocale(t.Locale, key, replacements...)
	return str
}

func (t Localizer) getLocalizationKey(locale string, key string) string {
	return fmt.Sprintf("%v.%v", locale, key)
}

func (t Localizer) replace(str string, replacements ...*Replacements) string {
	b := &bytes.Buffer{}
	tmpl, err := template.New("").Parse(str)
	if err != nil {
		return str
	}

	replacementsMerge := Replacements{}
	for _, replacement := range replacements {
		for k, v := range *replacement {
			replacementsMerge[k] = v
		}
	}

	err = template.Must(tmpl, err).Execute(b, replacementsMerge)
	if err != nil {
		return str
	}
	buff := b.String()
	return buff
}

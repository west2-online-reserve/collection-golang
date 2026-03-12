package main

type BiliData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Cursor struct {
			IsBegin         bool `json:"is_begin"`
			Prev            int  `json:"prev"`
			Next            int  `json:"next"`
			IsEnd           bool `json:"is_end"`
			PaginationReply struct {
				NextOffset string `json:"next_offset"`
			} `json:"pagination_reply"`
			SessionID   string `json:"session_id"`
			Mode        int    `json:"mode"`
			ModeText    string `json:"mode_text"`
			AllCount    int    `json:"all_count"`
			SupportMode []int  `json:"support_mode"`
			Name        string `json:"name"`
		} `json:"cursor"`
		Replies []struct {
			Rpid      int64  `json:"rpid"`
			Oid       int    `json:"oid"`
			Type      int    `json:"type"`
			Mid       int    `json:"mid"`
			Root      int    `json:"root"`
			Parent    int    `json:"parent"`
			Dialog    int    `json:"dialog"`
			Count     int    `json:"count"`
			Rcount    int    `json:"rcount"`
			State     int    `json:"state"`
			Fansgrade int    `json:"fansgrade"`
			Attr      int    `json:"attr"`
			Ctime     int    `json:"ctime"`
			MidStr    string `json:"mid_str"`
			OidStr    string `json:"oid_str"`
			RpidStr   string `json:"rpid_str"`
			RootStr   string `json:"root_str"`
			ParentStr string `json:"parent_str"`
			DialogStr string `json:"dialog_str"`
			Like      int    `json:"like"`
			Action    int    `json:"action"`
			Member    struct {
				Mid            string `json:"mid"`
				Uname          string `json:"uname"`
				Sex            string `json:"sex"`
				Sign           string `json:"sign"`
				Avatar         string `json:"avatar"`
				Rank           string `json:"rank"`
				FaceNftNew     int    `json:"face_nft_new"`
				IsSeniorMember int    `json:"is_senior_member"`
				Senior         struct {
				} `json:"senior"`
				LevelInfo struct {
					CurrentLevel int `json:"current_level"`
					CurrentMin   int `json:"current_min"`
					CurrentExp   int `json:"current_exp"`
					NextExp      int `json:"next_exp"`
				} `json:"level_info"`
				Pendant struct {
					Pid               int    `json:"pid"`
					Name              string `json:"name"`
					Image             string `json:"image"`
					Expire            int    `json:"expire"`
					ImageEnhance      string `json:"image_enhance"`
					ImageEnhanceFrame string `json:"image_enhance_frame"`
					NPid              int    `json:"n_pid"`
				} `json:"pendant"`
				Nameplate struct {
					Nid        int    `json:"nid"`
					Name       string `json:"name"`
					Image      string `json:"image"`
					ImageSmall string `json:"image_small"`
					Level      string `json:"level"`
					Condition  string `json:"condition"`
				} `json:"nameplate"`
				OfficialVerify struct {
					Type int    `json:"type"`
					Desc string `json:"desc"`
				} `json:"official_verify"`
				Vip struct {
					VipType       int    `json:"vipType"`
					VipDueDate    int64  `json:"vipDueDate"`
					DueRemark     string `json:"dueRemark"`
					AccessStatus  int    `json:"accessStatus"`
					VipStatus     int    `json:"vipStatus"`
					VipStatusWarn string `json:"vipStatusWarn"`
					ThemeType     int    `json:"themeType"`
					Label         struct {
						Path                  string      `json:"path"`
						Text                  string      `json:"text"`
						LabelTheme            string      `json:"label_theme"`
						TextColor             string      `json:"text_color"`
						BgStyle               int         `json:"bg_style"`
						BgColor               string      `json:"bg_color"`
						BorderColor           string      `json:"border_color"`
						UseImgLabel           bool        `json:"use_img_label"`
						ImgLabelURIHans       string      `json:"img_label_uri_hans"`
						ImgLabelURIHant       string      `json:"img_label_uri_hant"`
						ImgLabelURIHansStatic string      `json:"img_label_uri_hans_static"`
						ImgLabelURIHantStatic string      `json:"img_label_uri_hant_static"`
						LabelID               int         `json:"label_id"`
						LabelGoto             interface{} `json:"label_goto"`
					} `json:"label"`
					AvatarSubscript int    `json:"avatar_subscript"`
					NicknameColor   string `json:"nickname_color"`
				} `json:"vip"`
				FansDetail  interface{} `json:"fans_detail"`
				UserSailing struct {
					Pendant         interface{} `json:"pendant"`
					Cardbg          interface{} `json:"cardbg"`
					CardbgWithFocus interface{} `json:"cardbg_with_focus"`
				} `json:"user_sailing"`
				UserSailingV2 struct {
				} `json:"user_sailing_v2"`
				IsContractor   bool        `json:"is_contractor"`
				ContractDesc   string      `json:"contract_desc"`
				NftInteraction interface{} `json:"nft_interaction"`
				AvatarItem     struct {
					ContainerSize struct {
						Width  float64 `json:"width"`
						Height float64 `json:"height"`
					} `json:"container_size"`
					FallbackLayers struct {
						Layers []struct {
							Visible     bool `json:"visible"`
							GeneralSpec struct {
								PosSpec struct {
									CoordinatePos int     `json:"coordinate_pos"`
									AxisX         float64 `json:"axis_x"`
									AxisY         float64 `json:"axis_y"`
								} `json:"pos_spec"`
								SizeSpec struct {
									Width  int `json:"width"`
									Height int `json:"height"`
								} `json:"size_spec"`
								RenderSpec struct {
									Opacity int `json:"opacity"`
								} `json:"render_spec"`
							} `json:"general_spec"`
							LayerConfig struct {
								Tags struct {
									AVATARLAYER struct {
									} `json:"AVATAR_LAYER"`
									GENERALCFG struct {
										ConfigType    int `json:"config_type"`
										GeneralConfig struct {
											WebCSSStyle struct {
												BorderRadius string `json:"borderRadius"`
											} `json:"web_css_style"`
										} `json:"general_config"`
									} `json:"GENERAL_CFG"`
								} `json:"tags"`
								IsCritical bool `json:"is_critical"`
							} `json:"layer_config"`
							Resource struct {
								ResType  int `json:"res_type"`
								ResImage struct {
									ImageSrc struct {
										SrcType     int `json:"src_type"`
										Placeholder int `json:"placeholder"`
										Remote      struct {
											URL      string `json:"url"`
											BfsStyle string `json:"bfs_style"`
										} `json:"remote"`
									} `json:"image_src"`
								} `json:"res_image"`
							} `json:"resource"`
						} `json:"layers"`
						IsCriticalGroup bool `json:"is_critical_group"`
					} `json:"fallback_layers"`
					Mid string `json:"mid"`
				} `json:"avatar_item"`
			} `json:"member"`
			Content struct {
				Message string        `json:"message"`
				Members []interface{} `json:"members"`
				Emote   struct {
					OK struct {
						ID        int    `json:"id"`
						PackageID int    `json:"package_id"`
						State     int    `json:"state"`
						Type      int    `json:"type"`
						Attr      int    `json:"attr"`
						Text      string `json:"text"`
						URL       string `json:"url"`
						Meta      struct {
							Size    int      `json:"size"`
							Suggest []string `json:"suggest"`
						} `json:"meta"`
						Mtime     int    `json:"mtime"`
						JumpTitle string `json:"jump_title"`
					} `json:"[OK]"`
					NAMING_FAILED struct {
						ID        int    `json:"id"`
						PackageID int    `json:"package_id"`
						State     int    `json:"state"`
						Type      int    `json:"type"`
						Attr      int    `json:"attr"`
						Text      string `json:"text"`
						URL       string `json:"url"`
						Meta      struct {
							Size    int      `json:"size"`
							Suggest []string `json:"suggest"`
						} `json:"meta"`
						Mtime     int    `json:"mtime"`
						JumpTitle string `json:"jump_title"`
					} `json:"[藏狐]"`
				} `json:"emote"`
				JumpURL struct {
					Irb1200 struct {
						Title          string `json:"title"`
						State          int    `json:"state"`
						PrefixIcon     string `json:"prefix_icon"`
						AppURLSchema   string `json:"app_url_schema"`
						AppName        string `json:"app_name"`
						AppPackageName string `json:"app_package_name"`
						ClickReport    string `json:"click_report"`
						IsHalfScreen   bool   `json:"is_half_screen"`
						ExposureReport string `json:"exposure_report"`
						Extra          struct {
							GoodsShowType       int    `json:"goods_show_type"`
							IsWordSearch        bool   `json:"is_word_search"`
							GoodsCmControl      int    `json:"goods_cm_control"`
							GoodsClickReport    string `json:"goods_click_report"`
							GoodsExposureReport string `json:"goods_exposure_report"`
						} `json:"extra"`
						Underline    bool   `json:"underline"`
						MatchOnce    bool   `json:"match_once"`
						PcURL        string `json:"pc_url"`
						IconPosition int    `json:"icon_position"`
					} `json:"irb1200"`
					NAMING_FAILED struct {
						Title          string `json:"title"`
						State          int    `json:"state"`
						PrefixIcon     string `json:"prefix_icon"`
						AppURLSchema   string `json:"app_url_schema"`
						AppName        string `json:"app_name"`
						AppPackageName string `json:"app_package_name"`
						ClickReport    string `json:"click_report"`
						IsHalfScreen   bool   `json:"is_half_screen"`
						ExposureReport string `json:"exposure_report"`
						Extra          struct {
							GoodsShowType       int    `json:"goods_show_type"`
							IsWordSearch        bool   `json:"is_word_search"`
							GoodsCmControl      int    `json:"goods_cm_control"`
							GoodsClickReport    string `json:"goods_click_report"`
							GoodsExposureReport string `json:"goods_exposure_report"`
						} `json:"extra"`
						Underline    bool   `json:"underline"`
						MatchOnce    bool   `json:"match_once"`
						PcURL        string `json:"pc_url"`
						IconPosition int    `json:"icon_position"`
					} `json:"发那科"`
				} `json:"jump_url"`
				MaxLine int `json:"max_line"`
			} `json:"content,omitempty"`
			Replies []struct {
				Rpid      int64  `json:"rpid"`
				Oid       int    `json:"oid"`
				Type      int    `json:"type"`
				Mid       int    `json:"mid"`
				Root      int64  `json:"root"`
				Parent    int64  `json:"parent"`
				Dialog    int64  `json:"dialog"`
				Count     int    `json:"count"`
				Rcount    int    `json:"rcount"`
				State     int    `json:"state"`
				Fansgrade int    `json:"fansgrade"`
				Attr      int    `json:"attr"`
				Ctime     int    `json:"ctime"`
				MidStr    string `json:"mid_str"`
				OidStr    string `json:"oid_str"`
				RpidStr   string `json:"rpid_str"`
				RootStr   string `json:"root_str"`
				ParentStr string `json:"parent_str"`
				DialogStr string `json:"dialog_str"`
				Like      int    `json:"like"`
				Action    int    `json:"action"`
				Member    struct {
					Mid            string `json:"mid"`
					Uname          string `json:"uname"`
					Sex            string `json:"sex"`
					Sign           string `json:"sign"`
					Avatar         string `json:"avatar"`
					Rank           string `json:"rank"`
					FaceNftNew     int    `json:"face_nft_new"`
					IsSeniorMember int    `json:"is_senior_member"`
					Senior         struct {
					} `json:"senior"`
					LevelInfo struct {
						CurrentLevel int `json:"current_level"`
						CurrentMin   int `json:"current_min"`
						CurrentExp   int `json:"current_exp"`
						NextExp      int `json:"next_exp"`
					} `json:"level_info"`
					Pendant struct {
						Pid               int    `json:"pid"`
						Name              string `json:"name"`
						Image             string `json:"image"`
						Expire            int    `json:"expire"`
						ImageEnhance      string `json:"image_enhance"`
						ImageEnhanceFrame string `json:"image_enhance_frame"`
						NPid              int    `json:"n_pid"`
					} `json:"pendant"`
					Nameplate struct {
						Nid        int    `json:"nid"`
						Name       string `json:"name"`
						Image      string `json:"image"`
						ImageSmall string `json:"image_small"`
						Level      string `json:"level"`
						Condition  string `json:"condition"`
					} `json:"nameplate"`
					OfficialVerify struct {
						Type int    `json:"type"`
						Desc string `json:"desc"`
					} `json:"official_verify"`
					Vip struct {
						VipType       int    `json:"vipType"`
						VipDueDate    int64  `json:"vipDueDate"`
						DueRemark     string `json:"dueRemark"`
						AccessStatus  int    `json:"accessStatus"`
						VipStatus     int    `json:"vipStatus"`
						VipStatusWarn string `json:"vipStatusWarn"`
						ThemeType     int    `json:"themeType"`
						Label         struct {
							Path                  string `json:"path"`
							Text                  string `json:"text"`
							LabelTheme            string `json:"label_theme"`
							TextColor             string `json:"text_color"`
							BgStyle               int    `json:"bg_style"`
							BgColor               string `json:"bg_color"`
							BorderColor           string `json:"border_color"`
							UseImgLabel           bool   `json:"use_img_label"`
							ImgLabelURIHans       string `json:"img_label_uri_hans"`
							ImgLabelURIHant       string `json:"img_label_uri_hant"`
							ImgLabelURIHansStatic string `json:"img_label_uri_hans_static"`
							ImgLabelURIHantStatic string `json:"img_label_uri_hant_static"`
							LabelID               int    `json:"label_id"`
							LabelGoto             struct {
								Mobile string `json:"mobile"`
								PcWeb  string `json:"pc_web"`
							} `json:"label_goto"`
						} `json:"label"`
						AvatarSubscript int    `json:"avatar_subscript"`
						NicknameColor   string `json:"nickname_color"`
					} `json:"vip"`
					FansDetail     interface{} `json:"fans_detail"`
					UserSailing    interface{} `json:"user_sailing"`
					IsContractor   bool        `json:"is_contractor"`
					ContractDesc   string      `json:"contract_desc"`
					NftInteraction interface{} `json:"nft_interaction"`
					AvatarItem     struct {
						ContainerSize struct {
							Width  float64 `json:"width"`
							Height float64 `json:"height"`
						} `json:"container_size"`
						FallbackLayers struct {
							Layers []struct {
								Visible     bool `json:"visible"`
								GeneralSpec struct {
									PosSpec struct {
										CoordinatePos int     `json:"coordinate_pos"`
										AxisX         float64 `json:"axis_x"`
										AxisY         float64 `json:"axis_y"`
									} `json:"pos_spec"`
									SizeSpec struct {
										Width  int `json:"width"`
										Height int `json:"height"`
									} `json:"size_spec"`
									RenderSpec struct {
										Opacity int `json:"opacity"`
									} `json:"render_spec"`
								} `json:"general_spec"`
								LayerConfig struct {
									Tags struct {
										AVATARLAYER struct {
										} `json:"AVATAR_LAYER"`
										GENERALCFG struct {
											ConfigType    int `json:"config_type"`
											GeneralConfig struct {
												WebCSSStyle struct {
													BorderRadius string `json:"borderRadius"`
												} `json:"web_css_style"`
											} `json:"general_config"`
										} `json:"GENERAL_CFG"`
									} `json:"tags"`
									IsCritical bool `json:"is_critical"`
								} `json:"layer_config,omitempty"`
								Resource struct {
									ResType  int `json:"res_type"`
									ResImage struct {
										ImageSrc struct {
											SrcType     int `json:"src_type"`
											Placeholder int `json:"placeholder"`
											Remote      struct {
												URL      string `json:"url"`
												BfsStyle string `json:"bfs_style"`
											} `json:"remote"`
										} `json:"image_src"`
									} `json:"res_image"`
								} `json:"resource"`
							} `json:"layers"`
							IsCriticalGroup bool `json:"is_critical_group"`
						} `json:"fallback_layers"`
						Mid string `json:"mid"`
					} `json:"avatar_item"`
				} `json:"member"`
				Content struct {
					Message string        `json:"message"`
					Members []interface{} `json:"members"`
					Emote   struct {
						NAMING_FAILED struct {
							ID        int    `json:"id"`
							PackageID int    `json:"package_id"`
							State     int    `json:"state"`
							Type      int    `json:"type"`
							Attr      int    `json:"attr"`
							Text      string `json:"text"`
							URL       string `json:"url"`
							Meta      struct {
								Size    int      `json:"size"`
								Suggest []string `json:"suggest"`
							} `json:"meta"`
							Mtime     int    `json:"mtime"`
							JumpTitle string `json:"jump_title"`
						} `json:"[嗑瓜子]"`
					} `json:"emote"`
					JumpURL struct {
					} `json:"jump_url"`
					MaxLine int `json:"max_line"`
				} `json:"content"`
				Replies  interface{} `json:"replies"`
				Assist   int         `json:"assist"`
				UpAction struct {
					Like  bool `json:"like"`
					Reply bool `json:"reply"`
				} `json:"up_action"`
				Invisible    bool `json:"invisible"`
				ReplyControl struct {
					MaxLine           int    `json:"max_line"`
					TimeDesc          string `json:"time_desc"`
					TranslationSwitch int    `json:"translation_switch"`
					SupportShare      bool   `json:"support_share"`
				} `json:"reply_control"`
				Folder struct {
					HasFolded bool   `json:"has_folded"`
					IsFolded  bool   `json:"is_folded"`
					Rule      string `json:"rule"`
				} `json:"folder"`
				DynamicIDStr string `json:"dynamic_id_str"`
				NoteCvidStr  string `json:"note_cvid_str"`
				TrackInfo    string `json:"track_info"`
			} `json:"replies"`
			Assist   int `json:"assist"`
			UpAction struct {
				Like  bool `json:"like"`
				Reply bool `json:"reply"`
			} `json:"up_action"`
			Invisible    bool `json:"invisible"`
			ReplyControl struct {
				MaxLine           int    `json:"max_line"`
				SubReplyEntryText string `json:"sub_reply_entry_text"`
				SubReplyTitleText string `json:"sub_reply_title_text"`
				TimeDesc          string `json:"time_desc"`
				TranslationSwitch int    `json:"translation_switch"`
				SupportShare      bool   `json:"support_share"`
			} `json:"reply_control,omitempty"`
			Folder struct {
				HasFolded bool   `json:"has_folded"`
				IsFolded  bool   `json:"is_folded"`
				Rule      string `json:"rule"`
			} `json:"folder"`
			DynamicIDStr string `json:"dynamic_id_str"`
			NoteCvidStr  string `json:"note_cvid_str"`
			TrackInfo    string `json:"track_info"`
			CardLabel    []struct {
				Rpid             int64  `json:"rpid"`
				TextContent      string `json:"text_content"`
				TextColorDay     string `json:"text_color_day"`
				TextColorNight   string `json:"text_color_night"`
				LabelColorDay    string `json:"label_color_day"`
				LabelColorNight  string `json:"label_color_night"`
				Image            string `json:"image"`
				Type             int    `json:"type"`
				Background       string `json:"background"`
				BackgroundWidth  int    `json:"background_width"`
				BackgroundHeight int    `json:"background_height"`
				JumpURL          string `json:"jump_url"`
				Effect           int    `json:"effect"`
				EffectStartTime  int    `json:"effect_start_time"`
			} `json:"card_label,omitempty"`

			Contenta struct {
				Message string        `json:"message"`
				Members []interface{} `json:"members"`
				JumpURL struct {
					NAMING_FAILED struct {
						Title          string `json:"title"`
						State          int    `json:"state"`
						PrefixIcon     string `json:"prefix_icon"`
						AppURLSchema   string `json:"app_url_schema"`
						AppName        string `json:"app_name"`
						AppPackageName string `json:"app_package_name"`
						ClickReport    string `json:"click_report"`
						IsHalfScreen   bool   `json:"is_half_screen"`
						ExposureReport string `json:"exposure_report"`
						Extra          struct {
							GoodsShowType       int    `json:"goods_show_type"`
							IsWordSearch        bool   `json:"is_word_search"`
							GoodsCmControl      int    `json:"goods_cm_control"`
							GoodsClickReport    string `json:"goods_click_report"`
							GoodsExposureReport string `json:"goods_exposure_report"`
						} `json:"extra"`
						Underline    bool   `json:"underline"`
						MatchOnce    bool   `json:"match_once"`
						PcURL        string `json:"pc_url"`
						IconPosition int    `json:"icon_position"`
					}
				} `json:"jump_url"`
				MaxLine int `json:"max_line"`
			} `json:"contentb,omitempty"`

			Contentb struct {
				Message string        `json:"message"`
				Members []interface{} `json:"members"`
				JumpURL struct {
					NAMING_FAILED struct {
						Title          string `json:"title"`
						State          int    `json:"state"`
						PrefixIcon     string `json:"prefix_icon"`
						AppURLSchema   string `json:"app_url_schema"`
						AppName        string `json:"app_name"`
						AppPackageName string `json:"app_package_name"`
						ClickReport    string `json:"click_report"`
						IsHalfScreen   bool   `json:"is_half_screen"`
						ExposureReport string `json:"exposure_report"`
						Extra          struct {
							GoodsShowType       int    `json:"goods_show_type"`
							IsWordSearch        bool   `json:"is_word_search"`
							GoodsCmControl      int    `json:"goods_cm_control"`
							GoodsClickReport    string `json:"goods_click_report"`
							GoodsExposureReport string `json:"goods_exposure_report"`
						} `json:"extra"`
						Underline    bool   `json:"underline"`
						MatchOnce    bool   `json:"match_once"`
						PcURL        string `json:"pc_url"`
						IconPosition int    `json:"icon_position"`
					} `json:"稚辉君"`
				} `json:"jump_url"`
				MaxLine int `json:"max_line"`
			} `json:"contenta,omitempty"`
		} `json:"replies"`
		Top struct {
			Admin interface{} `json:"admin"`
			Upper interface{} `json:"upper"`
			Vote  interface{} `json:"vote"`
		} `json:"top"`
		TopReplies []interface{} `json:"top_replies"`
		Effects    struct {
			Preloading string `json:"preloading"`
		} `json:"effects"`
		Assist    int `json:"assist"`
		Blacklist int `json:"blacklist"`
		Vote      int `json:"vote"`
		Config    struct {
			Showtopic  int  `json:"showtopic"`
			ShowUpFlag bool `json:"show_up_flag"`
			ReadOnly   bool `json:"read_only"`
		} `json:"config"`
		Upper struct {
			Mid int `json:"mid"`
		} `json:"upper"`
		Control struct {
			InputDisable           bool        `json:"input_disable"`
			RootInputText          string      `json:"root_input_text"`
			ChildInputText         string      `json:"child_input_text"`
			GiveupInputText        string      `json:"giveup_input_text"`
			ScreenshotIconState    int         `json:"screenshot_icon_state"`
			UploadPictureIconState int         `json:"upload_picture_icon_state"`
			AnswerGuideText        string      `json:"answer_guide_text"`
			AnswerGuideIconURL     string      `json:"answer_guide_icon_url"`
			AnswerGuideIosURL      string      `json:"answer_guide_ios_url"`
			AnswerGuideAndroidURL  string      `json:"answer_guide_android_url"`
			BgText                 string      `json:"bg_text"`
			EmptyPage              interface{} `json:"empty_page"`
			ShowType               int         `json:"show_type"`
			ShowText               string      `json:"show_text"`
			WebSelection           bool        `json:"web_selection"`
			DisableJumpEmote       bool        `json:"disable_jump_emote"`
			EnableCharged          bool        `json:"enable_charged"`
			EnableCmBizHelper      bool        `json:"enable_cm_biz_helper"`
			PreloadResources       interface{} `json:"preload_resources"`
		} `json:"control"`
		Note   int `json:"note"`
		CmInfo struct {
			Ads struct {
				Num4765 []struct {
					ID           int    `json:"id"`
					ContractID   string `json:"contract_id"`
					PosNum       int    `json:"pos_num"`
					Name         string `json:"name"`
					Pic          string `json:"pic"`
					Litpic       string `json:"litpic"`
					URL          string `json:"url"`
					Style        int    `json:"style"`
					Agency       string `json:"agency"`
					Label        string `json:"label"`
					Intro        string `json:"intro"`
					CreativeType int    `json:"creative_type"`
					RequestID    string `json:"request_id"`
					SrcID        int    `json:"src_id"`
					Area         int    `json:"area"`
					IsAdLoc      bool   `json:"is_ad_loc"`
					AdCb         string `json:"ad_cb"`
					Title        string `json:"title"`
					ServerType   int    `json:"server_type"`
					CmMark       int    `json:"cm_mark"`
					Stime        int    `json:"stime"`
					Mid          string `json:"mid"`
					ActivityType int    `json:"activity_type"`
					Epid         int    `json:"epid"`
					SubTitle     string `json:"sub_title"`
					AdDesc       string `json:"ad_desc"`
					AdverName    string `json:"adver_name"`
					NullFrame    bool   `json:"null_frame"`
					PicMainColor string `json:"pic_main_color"`
				} `json:"4765"`
			} `json:"ads"`
		} `json:"cm_info"`
		Callbacks struct {
		} `json:"callbacks"`
	} `json:"data"`
}

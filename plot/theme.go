package plot

import (
	"encoding/json"
)

// ColorTheme contains the theme colors for plot
type ColorTheme struct {
	Name   string   `json:"name"`
	Colors []string `json:"colors"`
}

// ThemeColors contains all theme colors of ligo
var ThemeColors = []ColorTheme{}

// GetThemeColors get ThemeColors
func GetThemeColors(name string) (v ColorTheme) {
	for _, v = range ThemeColors {
		if v.Name == name {
			return v
		}
	}
	return v
}

var plotThemeData = []byte(`[
	{
			"name": "adobe_color_cc_1",
			"colors": ["#FFE350", "#E8740C", "#FF0000", "#9C0CE8", "#0D43FF", "#A6B212", "#1991FF", "#ECFF00", "#CC1E14", "#B25C58"]
	},
	{
			"name": "ball_subtype_colors",
			"colors": ["#0071bc", "#82a6d7", "#003b64", "#e8c122", "#8b770e", "#868686", "#b2b1b3", "#005083", "#c05f56", "#d79f93"]
	},
	{
			"name": "default",
			"colors": ["#0073c3", "#efc000", "#696969", "#ce534c", "#7ba6db", "#035892", "#052135", "#666633", "#660000", "#990000"]
	},
	{
			"name": "ggsci_aaas_default",
			"colors": ["#3B4992", "#EE0000", "#008B45", "#631879", "#008280", "#BB0021", "#5F559B", "#A20056", "#808180", "#1B1919"]
	},
	{
			"name": "ggsci_d3_category10",
			"colors": ["#1F77B4", "#FF7F0E", "#2CA02C", "#D62728", "#9467BD", "#8C564B", "#E377C2", "#7F7F7F", "#BCBD22", "#17BECF"]
	},
	{
			"name": "ggsci_d3_category20",
			"colors": ["#1F77B4", "#FF7F0E", "#2CA02C", "#D62728", "#9467BD", "#8C564B", "#E377C2", "#7F7F7F", "#BCBD22", "#17BECF", "#AEC7E8", "#FFBB78", "#98DF8A", "#FF9896", "#C5B0D5", "#C49C94", "#F7B6D2", "#C7C7C7", "#DBDB8D", "#9EDAE5"]
	},
	{
			"name": "ggsci_d3_category20b",
			"colors": ["#393B79", "#637939", "#8C6D31", "#843C39", "#7B4173", "#5254A3", "#8CA252", "#BD9E39", "#AD494A", "#A55194", "#6B6ECF", "#B5CF6B", "#E7BA52", "#D6616B", "#CE6DBD", "#9C9EDE", "#CEDB9C", "#E7CB94", "#E7969C", "#DE9ED6"]
	},
	{
			"name": "ggsci_d3_category20c",
			"colors": ["#3182BD", "#E6550D", "#31A354", "#756BB1", "#636363", "#6BAED6", "#FD8D3C", "#74C476", "#9E9AC8", "#969696", "#9ECAE1", "#FDAE6B", "#A1D99B", "#BCBDDC", "#BDBDBD", "#C6DBEF", "#FDD0A2", "#C7E9C0", "#DADAEB", "#D9D9D9"]
	},
	{
			"name": "ggsci_futurama_planetexpress",
			"colors": ["#FF6F00", "#C71000", "#008EA0", "#8A4198", "#5A9599", "#FF6348", "#84D7E1", "#FF95A8", "#3D3B25", "#ADE2D0", "#1A5354", "#3F4041"]
	},
	{
			"name": "ggsci_gsea_default",
			"colors": ["#4500AD", "#2700D1", "#6B58EF", "#8888FF", "#C7C1FF", "#D5D5FF", "#FFC0E5", "#FF8989", "#FF7080", "#FF5A5A", "#EF4040", "#D60C00"]
	},
	{
			"name": "ggsci_igv_alternating",
			"colors": ["#5773CC", "#FFB900"]
	},
	{
			"name": "ggsci_igv_default",
			"colors": ["#5050FF", "#CE3D32", "#749B58", "#F0E685", "#466983", "#BA6338", "#5DB1DD", "#802268", "#6BD76B", "#D595A7", "#924822", "#837B8D", "#C75127", "#D58F5C", "#7A65A5", "#E4AF69", "#3B1B53", "#CDDEB7", "#612A79", "#AE1F63", "#E7C76F", "#5A655E", "#CC9900", "#99CC00", "#A9A9A9", "#CC9900", "#99CC00", "#33CC00", "#00CC33", "#00CC99", "#0099CC", "#0A47FF", "#4775FF", "#FFC20A", "#FFD147", "#990033", "#991A00", "#996600", "#809900", "#339900", "#00991A", "#009966", "#008099", "#003399", "#1A0099", "#660099",
"#990080", "#D60047", "#FF1463", "#00D68F", "#14FFB1"]
	},
	{
			"name": "ggsci_jama_defalut",
			"colors": ["#374E55", "#DF8F44", "#00A1D5", "#B24745", "#79AF97", "#6A6599", "#80796B"]
	},
	{
			"name": "ggsci_jco_default",
			"colors": ["#0073C2", "#EFC000", "#868686", "#CD534C", "#7AA6DC", "#003C67", "#8F7700", "#3B3B3B", "#A73030", "#4A6990"]
	},
	{
			"name": "ggsci_lancet_lanonc",
			"colors": ["#00468B", "#ED0000", "#42B540", "#0099B4", "#925E9F", "#FDAF91", "#AD002A", "#ADB6B6", "#1B1919"]
	},
	{
			"name": "ggsci_locuszoom",
			"colors": ["#D43F3A", "#EEA236", "#5CB85C", "#46B8DA", "#357EBD", "#9632B8", "#B8B8B8"]
	},
	{
			"name": "ggsci_nejm_default",
			"colors": ["#BC3C29", "#0072B5", "#E18727", "#20854E", "#7876B1", "#6F99AD", "#FFDC91", "#EE4C97"]
	},
	{
			"name": "ggsci_npg_nrc",
			"colors": ["#E64B35", "#4DBBD5", "#00A087", "#3C5488", "#F39B7F", "#8491B4", "#91D1C2", "#DC0000", "#7E6148", "#B09C85"]
	},
	{
			"name": "ggsci_rickandmorty_schwifty",
			"colors": ["#FAFD7C", "#82491E", "#24325F", "#B7E4F9", "#FB6467", "#526E2D", "#E762D7", "#E89242", "#FAE48B", "#A6EEE6", "#917C5D", "#69C8EC"]
	},
	{
			"name": "ggsci_simpsons_springfield",
			"colors": ["#FED439", "#709AE1", "#8A9197", "#D2AF81", "#FD7446", "#D5E4A2", "#197EC0", "#F05C3B", "#46732E", "#71D0F5", "#370335", "#075149", "#C80813", "#91331F", "#1A9993", "#FD8CC1"]
	},
	{
			"name": "ggsci_startrek_uniform",
			"colors": ["#CC0C00", "#5C88DA", "#84BD00", "#FFCD00", "#7C878E", "#00B5E2", "#00AF66"]
	},
	{
			"name": "ggsci_tron_legacy",
			"colors": ["#FF410D", "#6EE2FF", "#F7C530", "#95CC5E", "#D0DFE6", "#F79D1E", "#748AA6"]
	},
	{
			"name": "ggsci_uchicago_dark",
			"colors": ["#800000", "#767676", "#CC8214", "#616530", "#0F425C", "#9A5324", "#642822", "#3E3E23", "#350E20"]
	},
	{
			"name": "ggsci_uchicago_default",
			"colors": ["#800000", "#767676", "#FFA319", "#8A9045", "#155F83", "#C16622", "#8F3931", "#58593F", "#350E20"]
	},
	{
			"name": "ggsci_uchicago_light",
			"colors": ["#800000", "#D6D6CE", "#FFB547", "#ADB17D", "#5B8FA8", "#D49464", "#B1746F", "#8A8B79", "#725663"]
	},
	{
			"name": "ggsci_ucscgb_default",
			"colors": ["#FF0000", "#FF9900", "#FFCC00", "#00FF00", "#6699FF", "#CC33FF", "#99991E", "#999999", "#FF00CC", "#CC0000", "#FFCCCC", "#FFFF00", "#CCFF00", "#358000", "#0000CC", "#99CCFF", "#00FFFF", "#CCFFFF", "#9900CC", "#CC99FF", "#996600", "#666600", "#666666", "#CCCCCC", "#79CC3D", "#CCCC99"]
	},
	{
			"name": "nature_brest_signatures",
			"colors": ["#3d1572", "#7d4594", "#e84286", "#f7c0ba", "#006230", "#199d47", "#91cd84", "#143080", "#3277b6", "#9584b9", "#dcd8e8", "#fac935"]
	},
	{
			"name": "ng_mutations",
			"colors": ["#609ec2", "#b56248", "#d0cb6c", "#9cb46f", "#5d9e71", "#b36170", "#aed5ec", "3d456f", "da8142"]
	},
	{
			"name": "nm_lines",
			"colors": ["#3b7ab5", "#e7211c", "#ff831d", "#2ee0d1", "#9c48a0", "#a7582c"]
	},
	{
			"name": "proteinpaint_chromHMM_state",
			"colors": ["#c0222c", "#f12424", "#ff00c7", "#d192fb", "#f9982f", "#fcc88e", "#fbf876", "#a6d67b", "#1fb855", "#007d37", "#00a99e", "#11aaec", "#186db9", "#3800f8", "#961a8b", "#47005f"]
	},
	{
			"name": "proteinpaint_domains",
			"colors": ["#a6d854", "#8dd3c7", "#fb8072", "#80b1d3", "#bebada", "#e5c494", "#fdb462", "#b3b3b3"]
	},
	{
			"name": "proteinpaint_mutations",
			"colors": ["#3987cc", "#ff7f0e", "#db3d3d", "#6633ff", "#bbbbbb", "#9467bd", "#998199", "#8c564b", "#819981", "#5781ff"]
	},
	{
			"name": "proteinpaint_significance",
			"colors": ["#aaaaaa", "#e99002", "#5bc0de", "#f04124", "#90c3d4", "#f04124", "#43ac6a"]
	},
	{
			"name": "red_blue",
			"colors": ["#c20b01", "#196abd"]
	}
]`)

func init() {
	json.Unmarshal(plotThemeData, &ThemeColors)
}

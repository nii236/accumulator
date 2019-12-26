import { createTheme, lightThemePrimitives, darkThemePrimitives } from "baseui"

// See https://github.com/uber-web/baseui/blob/master/src/themes/creator.js for full list of theme properties

// export const primaryOrange = "#FE7B6B"
// export const secondaryPurple = "#6C44A0"
// export const darkPurple = "#3E046A"

export const LightTheme = createTheme(
	{
		...lightThemePrimitives,
		// add all the properties here you'd like to override from the light theme primitives
		// primaryFontFamily: '"Comic Sans MS", cursive, sans-serif',
	},
	{
		// add all the theme overrides here - under the hood it uses deep merge
		// animation: {
		// 	timing100: '0.50s',
		// },
		// colors: {
		// 	progressStepsCompletedFill: primaryOrange,
		// 	buttonPrimaryFill: primaryOrange,
		// 	buttonPrimaryHover: "#FF9585",
		// 	buttonSecondaryFill: secondaryPurple,
		// 	buttonSecondaryHover: "#865EBA",
		// 	buttonSecondaryText: "#FFFFFF",
		// },
	},
)

export const DarkTheme = createTheme(
	{
		...darkThemePrimitives,
		// add all the properties here you'd like to override from the light theme primitives
		// primaryFontFamily: '"Comic Sans MS", cursive, sans-serif',
	},
	{
		// add all the theme overrides here - under the hood it uses deep merge
		// animation: {
		// 	timing100: '0.50s',
		// },
	},
)

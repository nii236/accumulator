import * as React from "react"
import MetaTags from "react-meta-tags"
import { BrowserRouter as Router, Route } from "react-router-dom"
import { Client as Styletron } from "styletron-engine-atomic"
import { Provider as StyletronProvider } from "styletron-react"
import { BaseProvider, useStyletron } from "baseui"
import { LightTheme, DarkTheme } from "./themeOverrides"
import { Teachers } from "./pages/Teachers"
import { Integrations } from "./pages/Integrations"
import { Friends } from "./pages/Friends"

const engine = new Styletron()
const Home = () => {
	return (
		<>
			<Integrations />
			<Friends />
			<Teachers />
		</>
	)
}
const Routes = () => {
	const [css, theme] = useStyletron()
	const routeStyle: string = css({
		width: "100%",
		minHeight: "100vh",
	})
	return (
		<div className={routeStyle}>
			<Router>
				{/* <Route path="/signin" component={SignIn} />
				<Route path="/signup" component={SignUp} /> */}
				<Route path="/friends" component={Friends} />
				<Route path="/teachers" component={Teachers} />
				<Route path="/integrations" component={Integrations} />
				{/* <Route path={"/verify/:code"} exact render={props => <EmailVerify code={props.match.params.code} />} /> */}
				{/* <Route path="/verify" exact component={EmailVerify} /> */}
			</Router>
		</div>
	)
}

const App = () => {
	const [darkTheme, setDarkTheme] = React.useState<boolean>(false)
	return (
		<StyletronProvider value={engine}>
			<BaseProvider theme={darkTheme ? DarkTheme : LightTheme}>
				<MetaTags>
					<title>Accumulator</title>
					<meta name="viewport" content="width=device-width, initial-scale=1.0" />
					<meta id="meta-description" name="description" content="Some description." />
					<meta id="og-title" property="og:title" content="MyApp" />
					<meta id="og-image" property="og:image" content="path/to/image.jpg" />
				</MetaTags>
				<Home />
			</BaseProvider>
		</StyletronProvider>
	)
}

export { App }

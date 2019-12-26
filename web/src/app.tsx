import * as React from "react"
import MetaTags from "react-meta-tags"

import { BrowserRouter as Router, Route, RouteComponentProps, Switch } from "react-router-dom"
import { Client as Styletron } from "styletron-engine-atomic"
import { Provider as StyletronProvider } from "styletron-react"
import { BaseProvider, useStyletron } from "baseui"
import { LightTheme, DarkTheme } from "./themeOverrides"
import { Teachers } from "./pages/Teachers"
import { Integrations } from "./pages/Integrations"
import { Friends } from "./pages/Friends"
import { Nav } from "./components/Nav"
import { Attendance } from "./pages/Attendance"

const engine = new Styletron()
interface Props extends RouteComponentProps {}
const Home = (props: Props) => {
	return (
		<>
			<Integrations {...props} />
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
				<Nav />
				<div>
					<Switch>
						<Route exact path="/" component={Home} />
						{/* <Route path="/signin" component={SignIn} />
				<Route path="/signup" component={SignUp} /> */}
						<Route exact path="/integrations/:integration_id/friends" component={Friends} />
						<Route exact path="/integrations/:integration_id/attendance" component={Attendance} />
						{/* <Route path={"/verify/:code"} exact render={props => <EmailVerify code={props.match.params.code} />} /> */}
						{/* <Route path="/verify" exact component={EmailVerify} /> */}
					</Switch>
				</div>
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
				<Routes />
			</BaseProvider>
		</StyletronProvider>
	)
}

export { App }

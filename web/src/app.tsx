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
import { SignIn } from "./pages/SignIn"

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
	const [validAuth, setValidAuth] = React.useState<boolean>(false)
	const routeStyle: string = css({
		width: "100%",
		minHeight: "100vh",
	})
	const authCheck = async () => {
		try {
			const res = await fetch("/api/auth/check")
			if (!res.ok) {
				const err = await res.text()
				throw new Error(err)
			}

			setValidAuth(true)
		} catch (err) {
			console.error(err)
			setValidAuth(false)
		}
	}
	React.useEffect(() => {
		authCheck()
	}, [])
	return (
		<div className={routeStyle}>
			{validAuth && (
				<Router>
					<Nav />
					<div>
						<Switch>
							<Route exact path="/" component={Home} />
							<Route exact path="/integrations/:integration_id/friends" component={Friends} />
							<Route exact path="/integrations/:integration_id/attendance/:teacher_id" component={Attendance} />
						</Switch>
					</div>
				</Router>
			)}
			{!validAuth && (
				<Router>
					<Nav />
					<div>
						<Switch>
							<Route path="/" component={SignIn} />
							{/* <Route path="/signup" component={SignUp} /> */}
						</Switch>
					</div>
				</Router>
			)}
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

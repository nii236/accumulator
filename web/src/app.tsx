import * as React from "react"
import MetaTags from "react-meta-tags"
import { BrowserRouter as Router, Route, RouteComponentProps, Switch, Redirect } from "react-router-dom"
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
import { SignUp } from "./pages/SignUp"
import { Spinner } from "baseui/spinner"
import { UI, useUI } from "./controllers/ui"
import { H1 } from "baseui/typography"
import { Users } from "./pages/Users"
import { APIDocumentation } from "./pages/apidoc"

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
	const ui = UI.useContainer()
	const [validAuth, setValidAuth] = React.useState<{ email: string; role: string } | null>(null)
	const routeStyle: string = css({
		width: "100%",
		minHeight: "100vh",
	})

	const authCheck = async () => {
		ui.startThinking()
		try {
			const res = await fetch("/api/auth/check")
			if (!res.ok) {
				const err = await res.text()
				throw new Error(err)
			}
			const data: { data: { email: string; role: string } } = await res.json()
			setValidAuth(data.data)
		} catch (err) {
			console.error(err)
			setValidAuth(null)
		}
		ui.stopThinking()
	}
	React.useEffect(() => {
		authCheck()
	}, [])

	if (ui.thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	return (
		<div className={routeStyle}>
			{validAuth && (
				<Router>
					<Nav email={validAuth.email} role={validAuth.role} />
					<div>
						<Switch>
							<Route exact path="/" component={Home} />
							<Route exact path="/documentation" component={APIDocumentation} />
							<Route exact path="/users" component={Users} />
							<Route exact path="/integrations/:integration_id/friends" component={Friends} />
							<Route exact path="/integrations/:integration_id/attendance/:teacher_id" component={Attendance} />
						</Switch>
					</div>
				</Router>
			)}
			{!validAuth && (
				<Router>
					<div>
						<H1 overrides={{ Block: { style: { textAlign: "center" } } }}>VRChat Accumulator System</H1>
						<Switch>
							<Route exact path="/" component={SignIn} />
							<Route exact path="/sign_up" component={SignUp} />
						</Switch>
					</div>
				</Router>
			)}
		</div>
	)
}
const engine = new Styletron()
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
				<UI.Provider>
					<Routes />
				</UI.Provider>
			</BaseProvider>
		</StyletronProvider>
	)
}

export { App }

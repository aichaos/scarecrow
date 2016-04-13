import React from "react";
import { browserHistory } from "react-router";

export default React.createClass({
	componentWillMount() {
		// If the app hasn't been configured yet, redirect them to setup.
		if (serverSettings.initialized === false) {
			browserHistory.push("/setup");
		}

		return null;
	},

	render() {
		return (
			<div>
				<h1>Home</h1>
			</div>
		);
	}
});

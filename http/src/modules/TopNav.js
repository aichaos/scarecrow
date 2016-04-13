import React from "react";
import { Link } from "react-router";

export default React.createClass({
	render() {
		return (
			<nav className="navbar navbar-inverse navbar-fixed-top">
				<div className="container-fluid">
					<div className="navbar-header">
						<button type="button" className="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
							<span className="sr-only">Toggle navigation</span>
							<span className="icon-bar"></span>
							<span className="icon-bar"></span>
							<span className="icon-bar"></span>
						</button>
						<Link to="/" className="navbar-brand">Scarecrow</Link>
					</div>
					<div id="navbar" className="navbar-collapse collapse">
						<ul className="nav navbar-nav navbar-right">
							<li><Link to="/">Dashboard</Link></li>
							<li><Link to="/setup">Setup</Link></li>
							<li><Link to="/settings">Settings</Link></li>
							<li><Link to="/help">Help</Link></li>
						</ul>
					</div>
				</div>
			</nav>
		);
	}
});

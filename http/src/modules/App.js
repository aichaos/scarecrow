import React from "react";
import { Link } from "react-router";

import TopNav from "./TopNav"
import LeftNav from "./LeftNav"

export default React.createClass({
	render() {
		return (
			<div>
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
							<TopNav/>
						</div>
					</div>
				</nav>

				<div className="container-fluid">
					<div className="row">
						<LeftNav/>
						<div className="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
							{this.props.children}
						</div>
					</div>
				</div>
			</div>
		);
	}
});

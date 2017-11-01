export default class Form extends React.PureComponent {
	static propTypes = {
		isWSReady: React.PropTypes.bool.isRequired,
		onSend: React.PropTypes.func.isRequired,
		update: React.PropTypes.object
	}

	static defaultProps = {
		isWSReady: false
	}

	constructor(props) {
		super(props)

		this.state = {
			input: '',
			showError: false
		}
		this.errorTimeout = null
		this.timeout = 3000
	}

	componentWillUpdate(nextProps, {showError}) {
		if (showError && !this.state.showError) {
			this.errorTimeout = setTimeout(() => {
				this.setState({showError: false})
			}, this.timeout)
		}
	}

	onSubmit = (e) => {
		e.preventDefault()

		const {input, showError} = this.state
		const {isWSReady} = this.props

		if (!input) return

		if (!isWSReady) {
			this.setState({showError: true})
		} else if (isWSReady && showError) {
			this.errorTimeout = null
			this.setState({showError: false})
		}

		if (!isWSReady) return

		this.props.onSend({
			search: input.toLowerCase()
		})
	}

	onTextInput = (e) => {
		this.setState({input: e.target.value})
	}

	renderConnectionError = () => {
		if (!this.state.showError) return null

		return (
			<div className="error">
				No connection to the server
			</div>
		)
	}

	render() {
		return (
			<div>
				<form className="form"
					onSubmit={this.onSubmit}>
					<input type="text"
						value={this.state.input}
						onChange={this.onTextInput} />

					<button type="submit">
						Submit
					</button>
				</form>

				{this.renderConnectionError()}
			</div>
		)
	}
}
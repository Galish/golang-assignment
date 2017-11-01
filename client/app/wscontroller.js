export default class WSController extends React.Component {
	static propTypes = {
		url: React.PropTypes.string
	}

	constructor(props, context) {
		super(props, context)
		this.state = {
			isOpened: false,
			update: null
		}
	}

	componentWillMount() {
		this.initWS()
	}

	componentDidUpdate(prevProps) {
		if (this.props.url && !prevProps.url) {
			this.initWS()
		}
	}

	componentWillUnmount() {
		this.onClose()
	}

	initWS = () => {
		this.ws = new WebSocket(this.props.url) //, 'echo-protocol')
		this.ws.onmessage = this.onMessage
		this.ws.onopen = this.onOpen
		this.ws.onerror = this.onError
	}

	onClose = () => {
		console.log(' > WS close')
		if (this.state.isOpened) {
			this.ws.close()
		}
	}

	onMessage = (e) => {
		// console.log(' > WS message', JSON.parse(e.data))
		this.setState({
			update: JSON.parse(e.data),
			timestamp: e.timeStamp
		})
	}

	onOpen = (e) => {
		// console.log(' > WS open')
		this.setState({isOpened: true})
		this.props.onReady && this.props.onReady()
	}

	onError = (e) => {
		console.log(' > WS error', e)
		// TODO: show global error message
		// this.setState({sOpened: true})
	}

	onSend = (message) => {
		// console.log(' > WS send', message)
		this.ws.send(JSON.stringify(message))
	}

	get update() {
		const {update, timestamp} = this.state
		return {
			...update,
			timestamp
		}
	}

	get isWSReady() {
		return !!this.ws && (this.ws.readyState === 1 || this.ws.readyState === 2)
	}

	render() {
		const {isWSReady, onClose, onMessage, onOpen, onError, onSend, update} = this
		const {isOpened} = this.state
		const {children} = this.props


		return React.cloneElement(
			children,
			{isWSReady, onClose, onMessage, onOpen, onError, onSend, update}
		)
	}
}

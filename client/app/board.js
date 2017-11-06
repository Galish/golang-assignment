export default class Board extends React.PureComponent {
	static propTypes = {
		isWSReady: React.PropTypes.bool.isRequired,
		update: React.PropTypes.object
	}

	static defaultProps = {
		isWSReady: false
	}

	constructor(props) {
		super(props)

		this.state = {
			messages: null
		}
	}

	componentWillReceiveProps({update}) {
		if (update) {
			if (update.result) {
				this.setState({
					messages: update.result
				})
			}
		}
	}

	get messages() {
		return (this.state.messages || []).sort(this.sortByDate)
	}

	sortByDate(itemA, itemB) {
		return itemA.date === itemB.date
			? 0
			: (itemA.date > itemB.date
				? -1
				: 1
			)
	}

	renderNotification = () => {
		const {messages} = this.state

		if (messages && !messages.length) {
			return (
				<div>
					No messages found
				</div>
			)
		}

		return null
	}

	renderResultsCount = () => {
		if (!this.messages.length) return null

		return (
			<p className="count">
				Messages found: <b>{this.messages.length}</b>
			</p>
		)
	}

	formatText = (html) => {
		const repl = '$br$'
		const arr = ['<br>', '<br \/>', '<br\/>', '<\/p>']

		return this.replaceArray(html, arr, repl)
			.replace(/<(?:.|\n)*?>/g, '')
			.split(repl).join('<br />')
	}

	replaceArray = (replaceString, find, repl) => {
		find.forEach(f => {
			const regex = new RegExp(f, 'g')
			replaceString = replaceString.replace(regex, repl)
		})

		return replaceString
  }

	render() {
		return (
			<div>
				{this.renderNotification()}

				{this.renderResultsCount()}

				{this.messages.map((message, i) =>
					<div className="message"
						key={i}>
						<div className="message__date">
							{moment(message.date).format('LLL')}
						</div>
						
						<div className="message__title">
							{message.title}
						</div>

						<div className="message__content"
							dangerouslySetInnerHTML={{
								__html: this.formatText(message.html)
							}} />

						<div className="message__author">
							{/*<img src={message.avatar} />*/}
							Author: <b>{message.author}</b>
						</div>

						<a className="message__link"
							href={message.link}
							target="_blank">
							{message.link}
						</a>
					</div>
				)}
			</div>
		)
	}
}

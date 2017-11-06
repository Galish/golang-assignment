import uuid  from 'node-uuid'
// import {replaceQuotes} from 'cm/utils'
// import {isFirefox} from 'cm/common/browser'
// import cloneDeep from 'lodash.clonedeep'

export function parseHTML(html) {
	let data = { corrected: false, tokens: [{ type: "unparsed", value: html }] }
	data = tokenizeTags(data)
	data = tokenizeLiterals(data)
	data = tokenizeGroups(data)
	data = tokenizeTermsAndOperators(data)

	data = parseGroups(data)
	data = parseOrOperator(data)

	data = generateOutput(data)

	return data
}

function tokenizeTags(input) {
	var corrected = input.corrected;
	var parsed = [];

	for (var t = 0; t < input.tokens.length; t++) {
		if (input.tokens[t].type !== "unparsed") {
			parsed.push(input.tokens[t]);
			continue;
		}

		var tokenValue = input.tokens[t].value;
		var offset = 0, lastOffset = 0;
		var value = "";

		//replace &nbsp; & <br> by whitespaces
		tokenValue = replaceQuotes(tokenValue).replace(/(&nbsp;|<br>)/g, " ")

		while ((offset = tokenValue.indexOf("<span", lastOffset)) >= 0) {
			value = tokenValue.substring(lastOffset, offset).trim();
			if (value.length > 0) {
				parsed.push({ type: "unparsed", value: value });
			}

			var endTag = "</span>";
			var endOffset = tokenValue.indexOf(endTag, offset + endTag.length);
			if (endOffset <= offset) {
				throw "invalid input";
			} else {
				value = tokenValue.substring(offset, endOffset + 1);
				parsed.push({ type: "tag", value: value });
			}

			lastOffset = endOffset + endTag.length;
		}

		value = tokenValue.substring(lastOffset, tokenValue.length).trim();
		if (value.length > 0) {
			parsed.push({ type: "unparsed", value: value });
		}
	}

	return { corrected: corrected, tokens: parsed };
}

function tokenizeLiterals(input) {
	var corrected = input.corrected;
	var parsed = [];

	for (var t = 0; t < input.tokens.length; t++) {
		if (input.tokens[t].type !== "unparsed") {
			parsed.push(input.tokens[t]);
			continue;
		}

		var tokenValue = input.tokens[t].value;
		var offset = 0, lastOffset = 0;
		var value = "";

		while ((offset = tokenValue.indexOf("\"", lastOffset)) >= 0) {
			value = tokenValue.substring(lastOffset, offset).trim();
			if (value.length > 0) {
				parsed.push({ type: "unparsed", value: value });
			}

			var endOffset = tokenValue.indexOf("\"", offset + 1);
			if (endOffset <= offset) {
				value = tokenValue.substring(offset + 1, tokenValue.length).trim();
				if (value.length > 0) {
					parsed.push({ type: "literal", value: value });
				}

				corrected = true;
				endOffset = tokenValue.length - 1;
			} else {
				value = tokenValue.substring(offset + 1, endOffset).trim();
				if (value.length > 0) {
					parsed.push({ type: "literal", value: value });
				} else {
					corrected = true;
				}
			}

			lastOffset = endOffset + 1;
		}

		value = tokenValue.substring(lastOffset, tokenValue.length).trim();
		if (value.length > 0) {
			parsed.push({ type: "unparsed", value: value });
		}
	}

	return { corrected: corrected, tokens: parsed };
}

function tokenizeGroups(input) {
	var corrected = input.corrected;
	var parsed = [];

	for (var t = 0; t < input.tokens.length; t++) {
		if (input.tokens[t].type !== "unparsed") {
			parsed.push(input.tokens[t]);
			continue;
		}

		var tokenValue = input.tokens[t].value;
		var value = "";

		for (var tv = 0; tv < tokenValue.length; tv++) {
			if (tokenValue.charAt(tv) === "(") {
				if (value.length > 0) {
					parsed.push({ type: "unparsed", value: value });
				}

				parsed.push({ type: "groupStart" });
				value = "";
			}
			else if (tokenValue.charAt(tv) === ")") {
				if (value.length > 0) {
					parsed.push({ type: "unparsed", value: value });
				}

				parsed.push({ type: "groupEnd" });
				value = "";
			}
			else {
				value += tokenValue.charAt(tv);
			}
		}

		if (value.length > 0) {
			parsed.push({ type: "unparsed", value: value });
		}
	}

	return { corrected: corrected, tokens: parsed };
}

function tokenizeTermsAndOperators(input) {
	var corrected = input.corrected;
	var parsed = [];

	for (var t = 0; t < input.tokens.length; t++) {
		if (input.tokens[t].type !== "unparsed") {
			parsed.push(input.tokens[t]);
			continue;
		}

		var tokenValue = input.tokens[t].value;
		var terms = tokenValue.match(/\S+/g) || [];

		for (var tt = 0; tt < terms.length; tt++) {
			var term = terms[tt];
			var addOrAfterTerm = false;

			if (term[0] === ',') {
				term = term.slice(1);
				parsed.push({ type: "or", children: [] });
			}
			else if (term[term.length - 1] === ',') {
				term = term.slice(0, -1);
				addOrAfterTerm = true;
			}

			if (!term || term === "<br>") {
				//ignore
			}
			else if (term === "|") {
				parsed.push({ type: "or", children: [] });
			}
			else if (term === '&' || term === "&amp;") {
				// note: AND is the default implicit operator and do not
				// interfer with OR because its has a higher precedence
			}
			else if (isHashtag(term)) {
				parsed.push({ type: "hashtag", value: term });
			}
			else if (isPattern(term)) {
				parsed.push({ type: "pattern", value: term });
			}
			else {
				parsed.push({ type: "term", value: term });
			}

			if (addOrAfterTerm) {
				parsed.push({ type: "or", children: [] });
			}
		}
	}

	return { corrected: corrected, tokens: parsed };
}

function parseGroups(input) {
	var corrected = input.corrected;
	var root = { type: "group", children: [], parent: null };
	var current = root;

	for (var t = 0; t < input.tokens.length; t++) {
		var token = input.tokens[t];

		if (token.type === "groupStart") {
			var group = { type: "group", children: [], parent: current };
			current.children.push(group);
			current = group;
		}
		else if (token.type === "groupEnd") {
			if (current.parent === null) {
				corrected = true;
				continue;
			}
			current = current.parent;
		}
		else {
			current.children.push(token);
		}
	}

	return { corrected: corrected, expression: root };
}

function parseOrOperator(input) {
	var corrected = input.corrected;
	var expression = input.expression;

	if ("children" in expression) {
		var children = expression.children || [];

		for (var c = 0; c < children.length;) {
			var result = parseOrOperator({ corrected: corrected, expression: children[c] });
			corrected = corrected || result.corrected;
			children[c] = result.expression;

			if (children[c].type === "group" && children[c].children.length === 0) {
				children.splice(c, 1);
			} else {
				c++;
			}
		}

		for (var c = 0; c < children.length;) {
			var child = children[c];
			var previousOr = null;

			if (child.type === "or") {
				var left = c;
				for (var cursor = c - 1; cursor >= 0; cursor--) {
					if (children[cursor].type === "or") {
						previousOr = children[cursor];
						break;
					}
					else {
						left = cursor;
					}
				}

				var right = c;
				for (var cursor = c + 1; cursor < children.length; cursor++) {
					if (children[cursor].type === "or") {
						break;
					}
					else {
						right = cursor;
					}
				}

				if (c === right) { //has no right branch
					corrected = true;
					children.splice(c, 1);
				}
				else if (previousOr !== null) { //fold
					var rightBranch = { type: "group", children: children.splice(c + 1, right - c) };
					if (rightBranch.children.length === 1) {
						rightBranch = rightBranch.children[0];
					}

					previousOr.children.push(rightBranch);
					children.splice(c, 1);
				}
				else if (c === left) { //no left branch
					corrected = true;
					children.splice(c, 1);
				}
				else {
					var rightBranch = { type: "group", children: children.splice(c + 1, right - c) };
					if (rightBranch.children.length === 1) {
						rightBranch = rightBranch.children[0];
					}

					var leftBranch = { type: "group", children: children.splice(left, c - left) };
					if (leftBranch.children.length === 1) {
						leftBranch = leftBranch.children[0];
					}

					child.children = [leftBranch, rightBranch];
					c = 1 + children.indexOf(child);
				}
			}
			else {
				c++;
			}
		}
	}

	return { corrected: corrected, expression: expression };
}

function generateOutput(input) {
	var corrected = input.corrected;
	var expression = input.expression;

	if (expression.type === "tag") {
		var tokenValue = expression.value;
		var tagCoords = (tokenValue.match(/(?:data-tag-coords=")([^\"]+)(?:")/) || [])[1];
		var tagType = (tokenValue.match(/(?:data-tag-type=")([^\"]+)(?:")/) || [])[1];
		var tagActive = (tokenValue.match(/(?:data-tag-active=")([^\"]+)(?:")/) || [])[1];
		var tagSource = (tokenValue.match(/(?:data-tag-source=")([^\"]+)(?:")/) || [])[1];
		var tagId = (tokenValue.match(/(?:data-tag-id=")([^\"]+)(?:")/) || [])[1] || "";
		var tagValue = (tokenValue.match(/(?:\>)([^\<]+)/) || [])[1] || "";

		// Delete '(removed)' from tag name
		if (tagActive === 'false') {
			tagValue = tagValue.split('(').slice(0, -1).join('(')
		}

		if (tagType === "buildingBlock" && tagId !== "") {
			expression = {
				'buildingBlock': {
					id: tagId,
					name: tagValue.trim()
				}
			}
			if (tagActive === 'false') expression.buildingBlock.removed = true
		}
		else if (tagType === "accountList" && tagId !== "") {
			expression = {
				'accountList': {
					id: tagId,
					name: tagValue.trim()
				}
			}
			if (tagActive === 'false') expression.accountList.removed = true
		}
		else if (tagType === "account" && tagId !== "") {
			expression = {account: {
				uri: tagId,
				type: tagSource,
				name: tagValue.trim()
			}}
		}
		else if (tagType === "location" && tagId !== "") {
			const coordinates = tagCoords.split('|').map((coord) => {
				const [lon, lat] = coord.split(':').map(Number)
				return {lon, lat}
			})
			expression = {"area": {type: "polygon", name: tagValue.trim(), coordinates}}
		}
		else {
			throw "invalid tag type or id ";
		}
	}
	else if (expression.type === "group") {
		var children = [];

		for (var c = 0; c < expression.children.length; c++) {
			var result = generateOutput({ corrected: corrected, expression: expression.children[c] });
			children.push(result.expression);
		}

		expression = { and: children };
	}
	else if (expression.type === "or") {
		var children = [];

		for (var c = 0; c < expression.children.length; c++) {
			var result = generateOutput({ corrected: corrected, expression: expression.children[c] });
			children.push(result.expression);
		}

		expression = { or: children };
	}
	else if (expression.type === "literal") {
		expression = {
			literal: {
				and: expression.value.match(/\S+/g).map((term) => ({term}))
			}
		};
	}
	else if (expression.type === "term") {
		expression = { term: expression.value };
	}
	else if (expression.type === "hashtag") {
		expression = { hashtag: expression.value };
	}
	else if (expression.type === "pattern") {
		expression = { pattern: expression.value };
	}

	return { corrected: corrected, expression: expression };
}

export function isHashtag(input) {
	return /^(#|\uFF03)([a-z0-9_\\u00c0-\\u00d6\\u00d8-\\u00f6\\u00f8-\\u00ff]*)$/i.test(input);
}

export function isPattern(input) {
	return /[\*\?]/.test(input)
}

export function convertToText(expression) {
	if(!expression)
		return "";

	if (expression.term) {
		return expression.term;

	} else if (expression.hashtag) {
		return expression.hashtag;

	} else if (expression.pattern) {
		return expression.pattern;

	} else if (expression.and) {
		return expression.and.map(convertToText).join(' ');

	} else if (expression.or) {
		return expression.or.map(convertToText).join(' ');

	} else if (expression.literal) {
		return `"${convertToText(expression.literal)}"`;

	} else if (expression.buildingBlock) {
		return expression.buildingBlock.name;

	} else if (expression.accountList) {
		return expression.accountList.name;

	} else if (expression.account) {
		return expression.account.name;

	} else if (expression.area) {
		return expression.area.name;

	} else if (expression.not) {
		return convertToText(expression.not);

	} else {
		return expression;
	}
}

export function constructTag(tag, type, color = 'blue2', isActive = true, language) {
	let location = ''
	if (type === 'area') {
		type = 'location'
		const coords = tag.coordinates.map(({lon, lat}) => [lon, lat].join(':')).join('|')
		location = `data-tag-coords="${coords}"`
	}
	let className = `query__${(type === 'account' ? `source--${tag.type}` : type).toLowerCase()}`
	className += ` tags-query__item--${color}`
	className += ` tags-query__item--${isActive && !tag.removed ? 'active' : 'inactive'}`
	if (isActive && tag.removed) className += ' tags-query__item--removed'

	const source = type === 'account' ? `data-tag-source="${tag.type}"` : ''
	const icon = (type === 'buildingBlock' || type === 'accountList')
		? `<i class="tags-query__icon icon icon--content ${tag.removed ? ' is-inactive ': ''}hint-preview"
			${!tag.removed ? `data-rh="${i18n.t('COMMON.SHOW_CONTENT')}"`: ''}
			data-tag-id="${tag.id}"
			data-tag-lang="${language}"
			data-tag-type="${type}"></i>`
		: ''
	const name = tag.removed
		? `${tag.name} (${i18n.t('COMMON.REMOVED')})`
		: tag.name

	return `&nbsp;<span id="${uuid.v4()}"
		${location}
		data-tag-id="${tag.id}"
		data-tag-type="${type}"
		data-tag-active="${!tag.removed}"
		${source}
		class="tags-query__item ${className}"
		unSelectable="on"
		contentEditable="false">&nbsp;${name}&nbsp;${icon}</span>&nbsp;`
}

export function convertToHTML({color, expression, isActive, hasBrackets, language}) {
	if (expression.term) {
		return expression.term;

	} else if (expression.hashtag) {
		return expression.hashtag;

	} else if (expression.pattern) {
		return expression.pattern;

	} else if (expression.and) {
		if (expression.and.length > 1) {
			hasBrackets = true;
		}

		return expression.and.map((expression) =>
			convertToHTML({expression, color, isActive, hasBrackets, language})).join(' ');

	} else if (expression.or) {
		const result = expression.or.map((expression) =>
			convertToHTML({expression, color, isActive, hasBrackets, language})).join(' | ');
		return hasBrackets ? `(${result})` : result;

	} else if (expression.literal) {
		return `"${convertToHTML({expression: expression.literal, color, isActive, hasBrackets, language})}"`;

	} else if (expression.buildingBlock) {
		return constructTag(expression.buildingBlock, 'buildingBlock', color, isActive, language)

	} else if (expression.accountList) {
		return constructTag(expression.accountList, 'accountList', color, isActive, language)

	} else if (expression.account) {
		return constructTag(Object.assign({}, expression.account, {id: expression.account.uri}), 'account', color, isActive, language)

	} else if (expression.area) {
		return constructTag(expression.area, 'area', color, isActive, language)

	} else {
		return '';
	}
}

function find(superset, subset, offset) {
	search: for (var s = offset; s < superset.length - subset.length + 1; s++) {
		for (var e = 0; e < subset.length; e++) {
			if (superset[s+e] !== subset[e]) {
				continue search;
			}
		}
		return s;
	}
	return -1;
}

function isValidUrlCharacter(character) {
	return !((character < 'A' || character > 'Z') && (character < 'a' || character > 'z') && (character < '0' || character > '9') && "-._~:/?#[]@!$&'()*+,;=".indexOf(character) < 0);
}

export function transformMessage({colors, matches, text}) {
	let tokens = [{ type: "unparsed", text: [...text], offset: 0 }];
	// matches = cloneDeep(matches).sort((a, b) =>
	matches = [...matches].sort((a, b) =>
		~~(b.filter || b.location) - ~~(a.filter || a.location))

	//remove overlapping matches
	let textMatches = [];
	for (let tm = 0; tm < matches.length; tm++) {
		const match = matches[tm];
		const hasActive = matches.some(m => m.start === match.start && m.end === match.end && m.active)
		const id = match.query || match.filter

		//http://stackoverflow.com/questions/3269434/whats-the-most-efficient-way-to-test-two-integer-ranges-for-overlap
		if ((match.active && hasActive ||!match.active && !hasActive) &&
			!textMatches.some(m => match.start < m.end && m.start < match.end) &&
			((id && colors[id]) || !id)
		) {
			textMatches.push(match)
		}
	}

	//resolve urls
	for (let t = 0; t < tokens.length; t++) {
		const token = tokens[t];

		if (token.type !== "unparsed") {
			continue;
		}

		const offsets = [
			find(token.text, [..."http://"], 0),
			find(token.text, [..."https://"], 0)
		].filter(s => s >= 0);

		if (offsets.length > 0) {
			let start = Math.min.apply(null, offsets);

			let end = start + 1;
			for (; isValidUrlCharacter(token.text[end]) && end < token.text.length; end++);

			tokens.splice(t, 1,
				{ type: "unparsed", text: token.text.slice(0, start), offset: token.offset },
				{ type: "url", text: token.text.slice(start, end), offset: token.offset + start },
				{ type: "unparsed", text: token.text.slice(end), offset: token.offset + end }
			);
		}
	}

	//resolve matches
	for (let tm = 0; tm < textMatches.length; tm++) {
		const match = textMatches[tm];

		for (let t = 0; t < tokens.length; t++) {
			const token = tokens[t];

			if (token.type !== "unparsed") {
				continue;
			}

			if (match.start >= token.offset && match.end <= token.offset + token.text.length) {
				const start = match.start - token.offset;
				const end = match.end - token.offset;

				tokens.splice(t, 1,
					{ type: "unparsed", text: token.text.slice(0, start), offset: token.offset },
					{ type: "match", text: token.text.slice(start, end), match: match, offset: token.offset + start },
					{ type: "unparsed", text: token.text.slice(end), offset: token.offset + end }
				);
			}
		}
	}

	//set extra properties
	for (let t = 0; t < tokens.length; t++) {
		const token = tokens[t];
		if (token.text) token.text = token.text.join('')

		switch (token.type) {
			case "unparsed":
				token.color = null;
				break;

			case "url":
				token.color = null;
				token.url = token.text;
				break;

			case "match":
				token.color = colors[token.match.filter || token.match.query || token.match.location || token.match.theme] || null;
				break;
		}
	}

	return tokens.filter(t => t.text.length > 0);
}

export function checkQueryValid(expression) {
	if (expression.and) {
		return expression.and.reduce((res, item) =>
			!checkQueryValid(item) && res ? false : res
		, true)

	} else if (expression.or) {
		return expression.or.reduce((res, item) =>
			!checkQueryValid(item) && res ? false : res
		, true)

	} else if (expression.buildingBlock) {
		return !expression.buildingBlock.removed

	} else if (expression.accountList) {
		return !expression.accountList.removed

	} else {
		return true
	}
}

export function replaceQuotes(string) {
	return replaceAll(string, ['«', '»', '“', '”', '″', '„'], '"')
}

export function replaceAll(string, search, replace) {
	search.map((item) => {
		string = string.split(item).join(replace)
	})
	return string.trim()
}

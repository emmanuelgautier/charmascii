package renderer

// emojiHeight is the number of rows in every emoji glyph definition.
const emojiHeight = 5

// emojiGlyphs maps an emoji base codepoint to a fixed emojiHeight-row ASCII art glyph.
var emojiGlyphs = map[rune][]string{
	// ── Faces ──────────────────────────────────────────────────────────────
	0x1F600: {" ___ ", "( o )", "| D |", " --- ", "     "},   // 😀 grinning
	0x1F601: {" ___ ", "(^ ^)", "| D |", " --- ", "     "},   // 😁 beaming
	0x1F602: {" ___ ", "(ToT)", "| D |", " --- ", " ~~~ "},   // 😂 tears of joy
	0x1F603: {" ___ ", "( o )", "| D |", " --- ", "     "},   // 😃 big eyes
	0x1F604: {" ___ ", "(^ ^)", "| D |", " --- ", "     "},   // 😄 smiling eyes
	0x1F605: {" ___ ", "(^ ^)", "| J |", " --- ", "   ' "},   // 😅 sweat
	0x1F606: {" ___ ", "(X X)", "| D |", " --- ", "     "},   // 😆 squinting
	0x1F923: {" ___ ", "(/XD)", "     ", " --- ", " ~~~ "},   // 🤣 floor laughing
	0x1F60A: {" ___ ", "(^ ^)", "| u |", " --- ", "     "},   // 😊 smiling
	0x1F607: {" _O_ ", "( o )", "| u |", " --- ", "     "},   // 😇 halo
	0x1F642: {" ___ ", "(. .)", "| J |", " --- ", "     "},   // 🙂 slightly smiling
	0x1F643: {" ___ ", "(\\ /)", "| n |", " --- ", "     "},  // 🙃 upside-down
	0x1F60D: {" ___ ", "(<><)", "| D |", " --- ", "     "},   // 😍 heart-eyes
	0x1F970: {" ___ ", "(ooo)", "| u |", " --- ", " <3  "},   // 🥰 hearts
	0x1F618: {" ___ ", "( o;)", "| J |", " --- ", "     "},   // 😘 blowing kiss
	0x1F61C: {" ___ ", "( o;)", "| P |", " --- ", "     "},   // 😜 winking tongue
	0x1F60E: {" ___ ", "([B])", "| _ |", " --- ", "     "},   // 😎 sunglasses
	0x1F913: {" ___ ", "([0])", "| _ |", " --- ", "     "},   // 🤓 nerd
	0x1F914: {" ___ ", "( o?)", "| / |", " --- ", "     "},   // 🤔 thinking
	0x1F634: {" ___ ", "(-_-)", "| z |", " --- ", "  zz "},   // 😴 sleeping
	0x1F631: {" ___ ", "( O )", "| O |", " --- ", "     "},   // 😱 screaming
	0x1F62D: {" ___ ", "(T T)", "| v |", " --- ", " ~~~ "},   // 😭 crying
	0x1F622: {" ___ ", "(. T)", "| _ |", " --- ", " ~   "},   // 😢 crying
	0x1F624: {" ___ ", "(><)", " | _ |", " --- ", "^  ^  "},  // 😤 steam
	0x1F620: {" ___ ", "(\\/)", "| _ |", " --- ", "     "},   // 😠 angry
	0x1F608: {"/\\-/\\", "( ^ )", "| J |", " --- ", "     "}, // 😈 devil
	0x1F47F: {"/\\-/\\", "(><)", "| _ |", " --- ", "     "},  // 👿 angry devil
	0x1F480: {" ___ ", "(x x)", "| ^ |", " |_| ", "  _  "},   // 💀 skull
	0x1F47B: {" ___ ", "( o )", "| u |", "|_  |", " ' ' "},   // 👻 ghost
	0x1F916: {"[___]", "|o o|", "|===|", "|___|", "     "},   // 🤖 robot
	0x1F921: {" ___ ", "( o )", "| O |", " --- ", "^^^^^"},   // 🤡 clown
	0x1F973: {" ___ ", "(^ ^)", "| D |", " --- ", "*   *"},   // 🥳 party face

	// ── Animals ────────────────────────────────────────────────────────────
	0x1F431: {" /\\ ", "(o o)", "=( )=", " ) ( ", "(___)"},  // 🐱 cat
	0x1F436: {" /~\\ ", "(o o)", "| _ |", " ) ( ", " --- "}, // 🐶 dog
	0x1F438: {" ___ ", "O . O", "|   |", " --- ", "     "},  // 🐸 frog

	// ── Symbols ────────────────────────────────────────────────────────────
	0x2764:  {"     ", " /\\/\\", " \\  /", "  \\/ ", "     "}, // ❤  heart
	0x1F525: {"  '  ", " /|  ", "( |) ", " \\|  ", "  ~  "},    // 🔥 fire
	0x2B50:  {"  *  ", " \\|/ ", "--*--", " /|\\ ", "  *  "},   // ⭐ star
	0x1F31F: {"  *  ", " *** ", "*****", " *** ", "  *  "},     // 🌟 glowing star
	0x2728:  {" * * ", "*   *", " * * ", "*   *", " * * "},     // ✨ sparkles
	0x1F4AF: {" 100 ", "=====", " 100 ", "!!!!!", "     "},     // 💯 hundred
	0x1F389: {" .   ", " (   ", "/ ** ", "**** ", "     "},     // 🎉 party popper
	0x1F680: {"  /\\ ", " /  \\", "| || |", " \\  /", " /\\ "}, // 🚀 rocket
	0x1F4AA: {"  _  ", " / \\ ", "( _ )", " \\_/ ", "     "},   // 💪 bicep
	0x1F44D: {" __  ", "( )  ", "| |  ", "|_|  ", "     "},     // 👍 thumbs up
	0x1F44F: {" | | ", "/   /", "| | |", "\\   \\", " ) ) "},   // 👏 clapping
	0x1F3C6: {" ___ ", "/   \\", "|   |", " | | ", " /_\\ "},   // 🏆 trophy
	0x26A1:  {" ___ ", "/ /  ", "/_/  ", "  \\  ", "  \\_ "},   // ⚡ lightning
	0x1F4A5: {" * * ", "* * *", " *** ", "* * *", " * * "},     // 💥 explosion
	0x274C:  {" X X ", "  X  ", " X X ", "     ", "     "},     // ❌ cross
	0x2705:  {"    /", "   / ", "  /  ", " /   ", "/    "},     // ✅ check
	0x1F3AF: {" ___ ", "/   \\", "| @ |", "\\___/", "     "},   // 🎯 target
	0x1F308: {"     ", ".::::.", ":::::::", " .::. ", "     "}, // 🌈 rainbow
}

// unknownEmojiGlyph returns a placeholder glyph for unrecognised emoji.
func unknownEmojiGlyph() []string {
	return []string{
		" ___ ",
		"|   |",
		"| ? |",
		"|___|",
		"     ",
	}
}

// isEmojiRune reports whether r is a non-ASCII codepoint that should be
// looked up in the emoji glyph map.
func isEmojiRune(r rune) bool {
	return r > 126
}

// isEmojiContinuation reports whether r is a combining/modifier codepoint
// that extends the preceding emoji (variation selector, ZWJ, skin tone).
func isEmojiContinuation(r rune) bool {
	return r == 0xFE0F || r == 0xFE0E || // variation selectors
		r == 0x200D || // zero-width joiner
		(r >= 0x1F3FB && r <= 0x1F3FF) // skin-tone modifiers
}

// getEmojiGlyph returns the ASCII art glyph for emoji rune r, padded (or
// truncated) to targetHeight rows so it aligns with FIGlet output.
func getEmojiGlyph(r rune, targetHeight int) []string {
	g, ok := emojiGlyphs[r]
	if !ok {
		g = unknownEmojiGlyph()
	}
	return padGlyphToHeight(g, targetHeight)
}

// padGlyphToHeight adjusts glyph to exactly targetHeight rows.
// Rows are added as empty strings at the top when shorter,
// or the top rows are dropped when taller.
func padGlyphToHeight(glyph []string, targetHeight int) []string {
	h := len(glyph)
	if h == targetHeight {
		return glyph
	}
	if h > targetHeight {
		return glyph[h-targetHeight:]
	}
	result := make([]string, targetHeight)
	copy(result[targetHeight-h:], glyph)
	return result
}

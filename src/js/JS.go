package js

var JS = "function utf8to16(str) {\n" +
	"    var out, i, len, c;\n" +
	"    var char2, char3;\n" +
	"    out = \"\";\n" +
	"    len = str.length;\n" +
	"    i = 0;\n" +
	"    while (i < len) {\n" +
	"        c = str.charCodeAt(i++);\n" +
	"        switch (c >> 4) {\n" +
	"            case 0:\n" +
	"            case 1:\n" +
	"            case 2:\n" +
	"            case 3:\n" +
	"            case 4:\n" +
	"            case 5:\n" +
	"            case 6:\n" +
	"            case 7:\n" +
	"                out += str.charAt(i - 1);\n" +
	"                break;\n" +
	"            case 12:\n" +
	"            case 13:\n" +
	"                char2 = str.charCodeAt(i++);\n" +
	"                out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));\n" +
	"                break;\n" +
	"            case 14:\n" +
	"                char2 = str.charCodeAt(i++);\n" +
	"                char3 = str.charCodeAt(i++);\n" +
	"                out += String.fromCharCode(((c & 0x0F) << 12) | ((char2 & 0x3F) << 6) | ((char3 & 0x3F) << 0));\n" +
	"                break\n" +
	"        }\n" +
	"    }\n" +
	"    return out\n" +
	"}\n" +
	"\n" +
	"var base64DecodeChars = new Array(-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1);\n" +
	"\n" +
	"function base64decode(str) {\n" +
	"    var c1, c2, c3, c4;\n" +
	"    var i, len, out;\n" +
	"    len = str.length;\n" +
	"    i = 0;\n" +
	"    out = \"\";\n" +
	"    while (i < len) {\n" +
	"        do {\n" +
	"            c1 = base64DecodeChars[str.charCodeAt(i++) & 0xff]\n" +
	"        } while (i < len && c1 == -1);\n" +
	"        if (c1 == -1) break;\n" +
	"        do {\n" +
	"            c2 = base64DecodeChars[str.charCodeAt(i++) & 0xff]\n" +
	"        } while (i < len && c2 == -1);\n" +
	"        if (c2 == -1) break;\n" +
	"        out += String.fromCharCode((c1 << 2) | ((c2 & 0x30) >> 4));\n" +
	"        do {\n" +
	"            c3 = str.charCodeAt(i++) & 0xff;\n" +
	"            if (c3 == 61) return out;\n" +
	"            c3 = base64DecodeChars[c3]\n" +
	"        } while (i < len && c3 == -1);\n" +
	"        if (c3 == -1) break;\n" +
	"        out += String.fromCharCode(((c2 & 0XF) << 4) | ((c3 & 0x3C) >> 2));\n" +
	"        do {\n" +
	"            c4 = str.charCodeAt(i++) & 0xff;\n" +
	"            if (c4 == 61) return out;\n" +
	"            c4 = base64DecodeChars[c4]\n" +
	"        } while (i < len && c4 == -1);\n" +
	"        if (c4 == -1) break;\n" +
	"        out += String.fromCharCode(((c3 & 0x03) << 6) | c4)\n" +
	"    }\n" +
	"    return out\n" +
	"}\n" +
	"\n" +
	"function strdecode(str) {\n" +
	"    return utf8to16(base64decode(str))\n" +
	"}\n" +
	"\n" +
	"function de(string, key) {\n" +
	"    string = strdecode(string);\n" +
	"    len = key.length;\n" +
	"    code = '';\n" +
	"    for (i = 0; i < string.length; i++) {\n" +
	"        k = i % len;\n" +
	"        code += String.fromCharCode(string.charCodeAt(i) ^ key.charCodeAt(k))\n" +
	"    }\n" +
	"    return strdecode(code)\n" +
	"}\n" +
	"\n" +
	"var kid = 2;\n" +
	"\n" +
	"function showopern(string, key) {\n" +
	"    var detxt = de(string, key);\n" +
	"    var c = detxt.split('|');\n" +
	"    var img_html = c[0];\n" +
	"    return img_html;\n" +
	"}\n"

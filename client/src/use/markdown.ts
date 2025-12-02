import MarkdownIt from "markdown-it";

const mdPromise = (async () => {
  const md = MarkdownIt({
    breaks: true
  });

  md.block.State.prototype.skipEmptyLines = function skipEmptyLines(from) {
    for (let max = this.lineMax; from < max; from++) {
      if (this.bMarks[from] + this.tShift[from] < this.eMarks[from]) {
        break;
      }
      this.push("hardbreak", "br", 0);
    }
    return from;
  };

  return md;
})();

export const render = async (text: string): Promise<string> => {
  return (await mdPromise).render(text);
};

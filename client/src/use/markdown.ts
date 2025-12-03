import MarkdownIt from "markdown-it";

const mdPromise = (async () => {
  const md = MarkdownIt({
    breaks: true
  });

  md.block.State.prototype.skipEmptyLines = function skipEmptyLines(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    from: any
  ) {
    let fromIndex = from;
    for (let max = this.lineMax; fromIndex < max; fromIndex++) {
      if (
        this.bMarks[fromIndex] + this.tShift[fromIndex] <
        this.eMarks[fromIndex]
      ) {
        break;
      }
      this.push("hardbreak", "br", 0);
    }
    return fromIndex;
  };

  return md;
})();

export const render = async (text: string): Promise<string> => {
  return (await mdPromise).render(text);
};

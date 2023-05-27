import { renderToString } from "react-dom/server";
import "./SearchResults.css";
import React, { useEffect, useState } from "react";
import { DocumentAbstract, SearchResult } from "../utils/types";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark as codeTheme } from "react-syntax-highlighter/dist/esm/styles/prism";
import { backend_url } from "../utils/config";
import {
  Table,
  TableHead,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Paper,
} from "@mui/material";
import { ThemeProvider, createTheme } from "@mui/material/styles";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

interface SearchResultProps {
  result: SearchResult;
  onFeedback: (resultId: string, score: number) => void;
}

const SearchResultItem: React.FC<SearchResultProps> = ({
  result,
  onFeedback,
}) => {
  const [content, setContent] = useState("");
  const [abstract, setAbstract] = useState(<div>Loading abstract...</div>);
  const [contentFetched, setContentFetched] = useState<boolean>(false);
  const [isFolded, setIsFolded] = useState(true); // new state
  const toggleFold = () => {
    setIsFolded(!isFolded);
  }; // function to toggle the fold
  // load

  useEffect(() => {
    // Abstract fetch
    fetch(`${backend_url}/extract_info?id=${result.id}`)
      .then((response) => response.json())
      .then((data) => {
        console.debug("data: ", data);
        let parsedData: DocumentAbstract = {
          entities: data.entities,
          hot_words: data.hot_words,
        };
        setAbstract(buildAbstract(parsedData));
      })
      .catch((error) => {
        console.error("Error loading abstract: ", error);
        setAbstract(<div>{"Failed: " + error}</div>);
      });

    // Content fetch
    if (!isFolded && !contentFetched) {
      fetch(`${backend_url}/document?id=${result.id}`)
        .then((response) => response.json())
        .then((data) => {
          let doc = renderCodeBlocks(data.content);
          setContent(doc.body.innerHTML);
          setContentFetched(true);
        });
    }
  }, [isFolded, content, contentFetched, result]);

  const buildAbstract = (data: DocumentAbstract) => {
    let entitiesTable = (
      <ThemeProvider theme={darkTheme}>
        <TableContainer component={Paper}>
          <h2 style={{ padding: "10px" }}>Entities</h2>
          <Table sx={{ tableLayout: "fixed" }}>
            <TableHead>
              <TableRow>
                <TableCell className="maxWidthCell">Entity</TableCell>
                <TableCell className="maxWidthCell">Frequency</TableCell>
              </TableRow>
            </TableHead>

            <TableBody>
              {Object.entries(data.entities).map(([key, value]) => {
                console.debug(key, value);
                return (
                  <TableRow key={key}>
                    <TableCell className="maxWidthCell">
                      <div className="maxWidthCell">{key}</div>
                    </TableCell>
                    <TableCell className="maxWidthCell">
                      <div className="maxWidthCell">{value}</div>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
      </ThemeProvider>
    );

    let hotWordsTable = (
      <ThemeProvider theme={darkTheme}>
        <TableContainer component={Paper}>
          <h2 style={{ padding: "10px" }}>Hot Words</h2>
          <Table sx={{ tableLayout: "fixed" }}>
            <TableHead>
              <TableRow>
                <TableCell>Word</TableCell>
                <TableCell>Frequency</TableCell>
              </TableRow>
            </TableHead>

            <TableBody>
              {Object.entries(data.hot_words).map(([key, value]) => (
                <TableRow key={key}>
                  <TableCell className="maxWidthCell">
                    <div className="maxWidthCell">{key}</div>
                  </TableCell>
                  <TableCell className="maxWidthCell">
                    <div className="maxWidthCell">{value}</div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </ThemeProvider>
    );

    return (
      <div>
        {entitiesTable}
        {hotWordsTable}
      </div>
    );
  };

  const renderCodeBlocks = (content: string) => {
    let parser = new DOMParser();
    let doc = parser.parseFromString(content, "text/html");
    let codeBlocks = doc.querySelectorAll("pre");

    codeBlocks.forEach((block) => {
      let codeString = block.textContent || "";
      let highlightedCodeString = renderToString(
        <SyntaxHighlighter language="cpp" style={codeTheme}>
          {codeString}
        </SyntaxHighlighter>
      );
      block.innerHTML = highlightedCodeString;
    });

    return doc;
  };

  const handleFeedback = (resultId: string, score: number) => {
    fetch(`${backend_url}/feedback`, {
      method: "POST",
      headers: {
        "Content-TYpe": "application/json",
      },
      body: JSON.stringify({ resultId, score }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        onFeedback(resultId, score);
      })
      .catch((error) => console.error("Error: ", error));
  };

  return (
    <div>
      <h2 onClick={toggleFold}>
        <div dangerouslySetInnerHTML={{ __html: result.title }} />
      </h2>
      {!isFolded && <div dangerouslySetInnerHTML={{ __html: content }} />}
      {abstract}
      <p>{result.url}</p>
      <p>{result.date}</p>
      <p>Relevance: {result.score}</p>
      <button className="feedback" onClick={() => handleFeedback(result.id, 1)}>
        👍
      </button>
      <button className="feedback" onClick={() => handleFeedback(result.id, 0)}>
        👎
      </button>
    </div>
  );
};

interface SearchResultsProps {
  results: SearchResult[];
}

const SearchResults: React.FC<SearchResultsProps> = ({ results }) => {
  const [adjustResults, setAdjustResults] = useState(results);

  useEffect(() => {
    setAdjustResults(results);
  }, [results]);

  const handleFeedback = (resultId: string, score: number) => {
    let newResults = [...adjustResults];

    let index = newResults.findIndex((result) => result.id === resultId);
    if (score === 1) {
      if (index > 0) {
        // Swap the item with the one before it to move it up
        [newResults[index], newResults[index - 1]] = [
          newResults[index - 1],
          newResults[index],
        ];
      }
    } else if (score === 0 && index < newResults.length - 1) {
      // Swap the item with the one after it to move it down
      [newResults[index], newResults[index + 1]] = [
        newResults[index + 1],
        newResults[index],
      ];
    }

    setAdjustResults(newResults);
  };

  console.debug(`Generate SearchResults for`);
  console.debug(results);

  return (
    <div>
      {adjustResults.map((result, index) => (
        <SearchResultItem
          key={`${result.id}-${index}`}
          result={result}
          onFeedback={handleFeedback}
        />
      ))}
    </div>
  );
};

export default SearchResults;

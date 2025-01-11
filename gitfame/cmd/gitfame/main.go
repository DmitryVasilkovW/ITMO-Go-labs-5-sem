//go:build !solution

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type IntSet map[string]int

type StringSet map[string]string

type UserDataSet map[string]UserData

type UserData struct {
	name    string
	commits IntSet
	files   int
	lines   int
}

type RepoFlags struct {
	repository   string
	revision     string
	orderBy      string
	useCommitter bool
	format       string
	extensions   []string
	languages    []string
	exclude      []string
	restrictTo   []string
}

func gitfame(request RepoFlags) ([]UserData, error) {
	var fileTree []string

	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", request.revision)
	cmd.Dir = request.repository
	fileNames, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	if len(fileNames) == 0 {
		fileTree = nil
	} else {
		fileTree = strings.Split(strings.TrimSpace(string(fileNames)), "\n")
	}

	userStats := make(UserDataSet)

	for _, fileName := range fileTree {
		if excludes(request, fileName) || restriced(request, fileName) ||
			extCheck(request, fileName) || langCheck(request, fileName) {
			continue
		}

		processedFile, err := process(fileName, request)
		if err != nil {
			return nil, err
		}

		for name, singleFile := range processedFile {
			userStat, ok := userStats[name]
			if !ok {
				userStat.commits = make(IntSet)
			}
			userStat.files += 1
			userStat.lines += singleFile.lines

			for commitHash := range singleFile.commits {
				userStat.commits[commitHash] = 1
			}

			userStats[name] = userStat
		}
	}

	var userInfo []UserData

	for name, stat := range userStats {
		var info UserData
		info.name = name
		info.files = stat.files
		info.lines = stat.lines
		info.commits = stat.commits

		userInfo = append(userInfo, info)
	}

	return userInfo, nil
}

func process(fileName string, info RepoFlags) (map[string]UserData, error) {

	cmd := exec.Command("git", "blame", fileName, "--porcelain", info.revision)
	cmd.Dir = info.repository
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	blame := string(cmdOutput)

	if len(blame) == 0 {
		cmd := exec.Command("git", "log", info.revision, "-1", "--pretty=format:%H %an", "--", fileName)
		cmd.Dir = info.repository
		cmdOutput, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		log := string(cmdOutput)

		lines := strings.Split(log, "\n")

		result := make(UserDataSet)

		for _, line := range lines {
			var singleFile UserData
			currentFields := strings.Fields(line)
			prefix := currentFields[0] + " "
			name := strings.TrimPrefix(line, prefix)
			singleFile.commits = make(IntSet)
			singleFile.commits[currentFields[0]] = 1
			result[name] = singleFile
		}

		return result, nil
	}

	lines := strings.Split(blame, "\n")
	var who string
	if info.useCommitter {
		who = "committer"
	} else {
		who = "author"
	}

	ccommitAmount := make(IntSet)
	commitsFromUser := make(StringSet)

	var currentHash string
	for _, line := range lines {
		currentFields := strings.Fields(line)

		if len(currentFields) == 0 {
			continue
		}

		if len(currentFields) == 4 {
			currentHash = currentFields[0]
			count, _ := strconv.Atoi(currentFields[3])
			ccommitAmount[currentHash] += count
			continue
		}
		if currentFields[0] == who {
			_, ok := commitsFromUser[currentHash]
			if !ok {
				prefix := who + " "
				commitsFromUser[currentHash] = strings.TrimPrefix(line, prefix)
			}
		}
	}

	result := make(UserDataSet)

	for commitHash, user := range commitsFromUser {
		stat, ok := result[user]
		if !ok {
			stat.commits = make(IntSet)
		}
		stat.commits[commitHash] = 1
		stat.lines += ccommitAmount[commitHash]

		result[user] = stat
	}

	return result, nil
}

func main() {
	var request RepoFlags

	flag.StringVar(&request.repository, "repository", ".", "путь до Git репозитория; по умолчанию текущая директория")
	flag.StringVar(&request.revision, "revision", "HEAD", "указатель на коммит; HEAD по умолчанию")
	flag.StringVar(&request.orderBy, "order-by", "lines", " ключ сортировки результатов; один из lines (дефолт), commits, files")
	flag.BoolVar(&request.useCommitter, "use-committer", false, "булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера")
	flag.StringVar(&request.format, "format", "tabular", "формат вывода; один из tabular (дефолт), csv, json, json-lines")
	flag.StringSliceVar(&request.extensions, "extensions", []string{}, "список расширений, сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например, '.go,.md'")
	flag.StringSliceVar(&request.languages, "languages", []string{}, "список языков (программирования, разметки и др.), сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например 'go,markdown'")
	flag.StringSliceVar(&request.exclude, "exclude", []string{}, "набор Glob паттернов, исключающих файлы из расчёта, например 'foo/*,bar/*'")
	flag.StringSliceVar(&request.restrictTo, "restrict-to", []string{}, "набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора")

	flag.Parse()

	answer, err := gitfame(request)
	if err != nil {
		fmt.Println("Can't get statistics: ", err)
		os.Exit(1)
	}

	switch request.orderBy {
	case "lines":
		sortByLine(answer)
	case "commits":
		sortByCommits(answer)
	case "files":
		sortByFiles(answer)
	default:
		os.Exit(1)
	}

	switch request.format {
	case "tabular":
		err := writeTabular(answer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing tabular result")
			os.Exit(1)
		}
	case "csv":
		err := writeCSV(answer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing csv result")
			os.Exit(1)
		}

	case "json":
		err := writeJSON(answer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing Json result")
			os.Exit(1)
		}
	case "json-lines":
		err := writeJSONLines(answer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing json-lines result")
			os.Exit(1)
		}
	default:
		os.Exit(1)
	}
}

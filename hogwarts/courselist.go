//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	visited := make(map[string]int)
	var courses []string

	for course := range prereqs {
		if isItCycleDependency(course, visited, prereqs, &courses) {
			panic("Cycle dependency detected")
		}
	}

	return courses
}

func isItCycleDependency(course string, visited map[string]int, prereqs map[string][]string, courses *[]string) bool {
	return isCourseNotVisited(visited, course) && !dfs(course, visited, prereqs, courses)
}

func isCourseNotVisited(visited map[string]int, course string) bool {
	return visited[course] == 0
}

func dfs(
	course string,
	visited map[string]int,
	prereqs map[string][]string,
	courses *[]string) bool {

	if isCourseInProcess(visited, course) {
		return false
	} else if isCourseVisited(visited, course) {
		return true
	}

	addCourseToProcess(visited, course)
	for _, prereq := range prereqs[course] {
		if !dfs(prereq, visited, prereqs, courses) {
			return false
		}
	}

	markCourseAsVisited(visited, course)
	*courses = append(*courses, course)
	return true
}

func addCourseToProcess(visited map[string]int, course string) {
	visited[course] = 1
}

func markCourseAsVisited(visited map[string]int, course string) {
	visited[course] = 2
}

func isCourseInProcess(visited map[string]int, course string) bool {
	return visited[course] == 1
}

func isCourseVisited(visited map[string]int, course string) bool {
	return visited[course] == 2
}
